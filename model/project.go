package model

import (
	"database/sql"
	"errors"

	"github.com/n-inja/go-blog/utils"
)

type Project struct {
	ID          string   `json:"id" form:"id"`
	Name        string   `json:"name" form:"name"`
	DisplayName string   `json:"displayName" form:"displayName"`
	UserID      string   `json:"userId" form:"userId"`
	Member      []string `json:"member" form:"member"`
	Description string   `json:"description" form:"description"`
	PostCount   int      `json:"postCount" form:"postCount"`
}

func (project *Project) Insert() error {
	if !utils.RegexProjectName.MatchString(project.Name) {
		return errors.New("project name := ^[a-zA-Z0-9_-]+$")
	}
	return utils.Transact(func(tx *sql.Tx) error {
		_, err := tx.Exec("insert into projects (id, name, display_name, user_id, description) value(?, ?, ?, ?, ?)", project.ID, project.Name, project.DisplayName, project.UserID, project.Description)
		if err != nil {
			return err
		}
		for _, userID := range project.Member {
			_, err = tx.Exec("insert into member (user_id, project_id) value(?, ?)", userID, project.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (project *Project) Delete() error {
	_, err := utils.DB.Exec("delete from projects where id = ?", project.ID)
	if err != nil {
		return err
	}
	_, err = utils.DB.Exec("update posts set is_deleted = true where project_id = ?", project.ID)
	if err != nil {
		return err
	}
	_, err = utils.DB.Exec("delete from member where project_id = ?", project.ID)
	if err != nil {
		return err
	}
	_, err = utils.DB.Exec("update comments set is_deleted = true where post_id in (select id from posts where project_id = ?)", project.ID)
	return err
}

func (project *Project) Update(invites, removes []string) error {
	if !utils.RegexProjectName.MatchString(project.Name) {
		return errors.New("project name := ^[a-zA-Z0-9_-]+$")
	}
	return utils.Transact(func(tx *sql.Tx) error {
		_, err := tx.Exec("update projects set name = ?, display_name = ?, user_id = ?, description = ? where id = ?", project.Name, project.DisplayName, project.UserID, project.Description, project.ID)
		if err != nil {
			return err
		}
		for _, userID := range invites {
			_, err = tx.Exec("insert into member (user_id, project_id) value(?, ?)", userID, project.ID)
			if err != nil {
				return err
			}
		}
		for _, userID := range removes {
			_, err = tx.Exec("delete from member where user_id = ? and project_id = ?", userID, project.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func GetProjects() ([]Project, error) {
	memberRows, err := utils.DB.Query("select project_id, user_id from member")
	if err != nil {
		return nil, err
	}
	defer memberRows.Close()
	userMap := map[string][]string{}
	for memberRows.Next() {
		var userID, projectID string
		memberRows.Scan(&projectID, &userID)
		if userMap[projectID] == nil {
			userMap[projectID] = make([]string, 0)
		}
		userMap[projectID] = append(userMap[projectID], userID)
	}

	projectRows, err := utils.DB.Query("select projects.id, name, display_name, projects.user_id, description, count(*) from projects left join posts on posts.project_id = projects.id group by projects.id")
	if err != nil {
		return nil, err
	}
	defer projectRows.Close()
	projects := make([]Project, 0)
	for projectRows.Next() {
		var project Project
		var description sql.NullString
		projectRows.Scan(&project.ID, &project.Name, &project.DisplayName, &project.UserID, &description, &project.PostCount)
		project.Description = ""
		if description.Valid {
			project.Description = description.String
		}
		if userMap[project.ID] == nil {
			project.Member = make([]string, 0)
		} else {
			project.Member = userMap[project.ID]
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func GetProject(ID string) (Project, error) {
	var project Project
	var description sql.NullString
	err := utils.DB.QueryRow("select p.id, p.name, p.display_name, p.user_id, p.description, count(*) from (select * from projects where id = ?) p left join posts on p.id = posts.project_id group by p.id", ID).Scan(&project.ID, &project.Name, &project.DisplayName, &project.UserID, &description, &project.PostCount)
	if err != nil {
		return Project{}, err
	}
	project.Description = ""
	if description.Valid {
		project.Description = description.String
	}

	rows, err := utils.DB.Query("select user_id from member where project_id = ?", project.ID)
	if err != nil {
		return Project{}, err
	}
	defer rows.Close()

	userIDs := make([]string, 0)
	for rows.Next() {
		var ID string
		rows.Scan(&ID)
		userIDs = append(userIDs, ID)
	}
	project.Member = userIDs

	return project, nil
}

func GetProjectByName(projectName string) (Project, error) {
	var projectID string
	err := utils.DB.QueryRow("select id from projects where name = ?", projectName).Scan(&projectID)
	if err != nil {
		return Project{}, err
	}

	var project Project
	err = utils.DB.QueryRow("select id, name, description from projects where name = ?", projectName).Scan(&project.ID, &project.Name, &project.Description)

	if err != nil {
		return Project{}, err
	}

	return project, nil
}
