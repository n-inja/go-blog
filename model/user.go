package model

import (
	"database/sql"
	"errors"

	"github.com/n-inja/blog/utils"
)

type User struct {
	ID          string   `json:"id" form:"id"`
	Name        string   `json:"name" form:"name"`
	Auth        string   `json:"auth" form:"auth"`
	ProjectIDs  []string `json:"projectIds" form:"projectIds"`
	Description string   `json:"description" form:"description"`
	IconSrc     string   `json:"iconSrc" form:"iconSrc"`
	TwitterId   string   `json:"twitterId" form:"twitterId"`
	GithubId    string   `json:"githubId" form:"githubId"`
}

func (user *User) Update() error {
	_, err := utils.DB.Exec("update profiles set description = ?, twitter_id = ?, github_id = ?, icon_src = ? where id = ?", user.Description, user.TwitterId, user.GithubId, user.IconSrc, user.ID)
	return err
}

func GetUsers() ([]User, error) {
	rows, err := utils.DB.Query("select user_id, project_id from member")
	if err != nil {
		return nil, err
	}

	memberMap := map[string][]string{}
	for rows.Next() {
		var userID, projectID string
		rows.Scan(&userID, &projectID)
		if memberMap[userID] == nil {
			memberMap[userID] = make([]string, 0)
		}
		memberMap[userID] = append(memberMap[userID], projectID)
	}
	rows.Close()

	rows, err = utils.DB.Query("select name, users.id, description, auth, icon_src, twitter_id, github_id from users left join profiles on users.id = profiles.id where auth = 'default' order by id desc")
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)
	for rows.Next() {
		var u User
		var description, iconSrc, twitterId, githubId sql.NullString
		rows.Scan(&u.Name, &u.ID, &description, &u.Auth, &iconSrc, &twitterId, &githubId)
		if memberMap[u.ID] == nil {
			u.ProjectIDs = make([]string, 0)
		} else {
			u.ProjectIDs = memberMap[u.ID]
		}
		u.Description = ""
		if description.Valid {
			u.Description = description.String
		}
		u.IconSrc = ""
		if iconSrc.Valid {
			u.IconSrc = iconSrc.String
		}
		u.TwitterId = ""
		if twitterId.Valid {
			u.TwitterId = twitterId.String
		}
		u.GithubId = ""
		if githubId.Valid {
			u.GithubId = githubId.String
		}
		users = append(users, u)
	}
	rows.Close()

	return users, nil
}

func GetUser(ID string) (User, error) {
	rows, err := utils.DB.Query("select name, users.id, description, auth, icon_src, twitter_id, github_id from users left join profiles on users.id = profiles.id where users.id = ? and auth = 'default'", ID)
	if err != nil {
		return User{}, err
	}
	if !rows.Next() {
		return User{}, errors.New("user not found")
	}
	var user User
	var description, iconSrc, twitterId, githubId sql.NullString
	rows.Scan(&user.Name, &user.ID, &description, &user.Auth, &iconSrc, &twitterId, &githubId)
	rows.Close()

	user.Description = ""
	if description.Valid {
		user.Description = description.String
	}
	user.IconSrc = ""
	if iconSrc.Valid {
		user.IconSrc = iconSrc.String
	}
	user.TwitterId = ""
	if twitterId.Valid {
		user.TwitterId = twitterId.String
	}
	user.GithubId = ""
	if githubId.Valid {
		user.GithubId = githubId.String
	}

	rows, err = utils.DB.Query("select project_id from member where user_id = ?", ID)
	if err != nil {
		return user, err
	}
	user.ProjectIDs = make([]string, 0)
	for rows.Next() {
		var projectID string
		rows.Scan(&projectID)
		user.ProjectIDs = append(user.ProjectIDs, projectID)
	}
	return user, nil
}
