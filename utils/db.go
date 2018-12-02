package utils

import (
	"database/sql"
	"errors"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var regexProjectName *regexp.Regexp

func Open(userName, password, address, databaseName string) error {
	var err error
	db, err = sql.Open("mysql", userName+":"+password+"@"+address+"/"+databaseName)
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(0)
	regexProjectName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	return initDB()
}

func Close() {
	db.Close()
}

func initDB() error {
	rows, err := db.Query("show tables like 'users'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		rows.Close()
		return errors.New("users table not found")
	}
	rows.Close()

	rows, err = db.Query("show tables like 'projects'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		_, err = db.Exec("create table projects (id varchar(20) NOT NULL PRIMARY KEY, name varchar(32) NOT NULL UNIQUE, display_name varchar(64) unicode NOT NULL, user_id varchar(32) NOT NULL, description text unicode NULL, foreign key(user_id) references users(id)) engine=innodb")
		if err != nil {
			return err
		}
	}
	rows.Close()

	rows, err = db.Query("show tables like 'member'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		_, err = db.Exec("create table member (user_id varchar(32) NOT NULL, project_id varchar(20) NOT NULL, PRIMARY KEY(user_id, project_id), index(user_id), index(project_id), foreign key(user_id) references users(id)) engine=innodb")
		if err != nil {
			return err
		}
	}
	rows.Close()

	rows, err = db.Query("show tables like 'posts'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		_, err = db.Exec("create table posts (id varchar(20) NOT NULL PRIMARY KEY, title text unicode NOT NULL, content longtext unicode NOT NULL, thumb_src varchar(64) NULL, user_id varchar(32) NOT NULL, created_at timestamp NOT NULL, updated_at timestamp NOT NULL, project_id varchar(20) NOT NULL, views int NOT NULL, is_deleted boolean NOT NULL, index(user_id), index(created_at), index(is_deleted), index(project_id), foreign key(user_id) references users(id)) engine=innodb")
		if err != nil {
			return err
		}
	}
	rows.Close()

	rows, err = db.Query("show tables like 'comments'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		_, err = db.Exec("create table comments (id varchar(20) NOT NULL PRIMARY KEY, content text unicode NOT NULL, user_id varchar(32) NOT NULL, post_id varchar(20) NOT NULL, created_at timestamp NOT NULL, is_deleted boolean NOT NULL, index(user_id), index(created_at), index(is_deleted), index(post_id), foreign key(post_id) references posts(id), foreign key(user_id) references users(id)) engine=innodb")
		if err != nil {
			return err
		}
	}
	rows.Close()

	rows, err = db.Query("show tables like 'profiles'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		_, err = db.Exec("create table profiles (id varchar(32) NOT NULL PRIMARY KEY, description text unicode NULL, twitter_id varchar(32) NULL, github_id varchar(64) NULL, icon_src varchar(64) NULL, foreign key(id) references users(id)) engine=innodb")
		if err != nil {
			return err
		}
	}
	rows.Close()

	return nil
}

type Post struct {
	ID         string `json:"id" form:"id"`
	Title      string `json:"title" form:"title"`
	Content    string `json:"content" form:"content"`
	ThumbSrc   string `json:"thumbSrc" form:"thumbSrc"`
	UserID     string `json:"userId" form:"userId"`
	CreatedAt  string `json:"createdAt" form:"createdAt"`
	UpdatedAt  string `json:"updatedAt" form:"updatedAt"`
	ProjectID  string `json:"projectId" form:"projectId"`
	Views      int    `json:"views" form:"views"`
	CommentNum int    `json:"commentNum" form:"commentNum"`
}

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

func GetUsers() ([]User, error) {
	rows, err := db.Query("select user_id, project_id from member")
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

	rows, err = db.Query("select name, users.id, description, auth, icon_src, twitter_id, github_id from users left join profiles on users.id = profiles.id where auth = 'default' order by id desc")
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
	rows, err := db.Query("select name, users.id, description, auth, icon_src, twitter_id, github_id from users left join profiles on users.id = profiles.id where users.id = ? and auth = 'default'", ID)
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

	rows, err = db.Query("select project_id from member where user_id = ?", ID)
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

func GetUserPosts(userID string, offset, limit int) ([]Post, error) {
	rows, err := db.Query("select posts.id, title, posts.content, thumb_src, posts.user_id, posts.created_at, updated_at, project_id, views, count(comments.id) from (select * from posts where user_id = ? and is_deleted = false order by created_at desc limit ?, ?) posts left join comments on comments.post_id = posts.id group by posts.id", userID, offset, limit)
	if err != nil {
		return make([]Post, 0), err
	}
	posts := make([]Post, 0)
	for rows.Next() {
		var post Post
		var thumbSrc sql.NullString
		rows.Scan(&post.ID, &post.Title, &post.Content, &thumbSrc, &post.UserID, &post.CreatedAt, &post.UpdatedAt, &post.ProjectID, &post.Views, &post.CommentNum)
		post.ThumbSrc = ""
		if thumbSrc.Valid {
			post.ThumbSrc = thumbSrc.String
		}
		posts = append(posts, post)
	}
	rows.Close()
	return posts, nil
}

type Project struct {
	ID          string   `json:"id" form:"id"`
	Name        string   `json:"name" form:"name"`
	DisplayName string   `json:"displayName" form:"displayName"`
	UserID      string   `json:"userId" form:"userId"`
	Member      []string `json:"member" form:"member"`
	Description string   `json:"description" form:"description"`
}

func GetProjects() ([]Project, error) {
	rows, err := db.Query("select project_id, user_id from member")
	if err != nil {
		return nil, err
	}
	userMap := map[string][]string{}
	for rows.Next() {
		var userID, projectID string
		rows.Scan(&projectID, &userID)
		if userMap[projectID] == nil {
			userMap[projectID] = make([]string, 0)
		}
		userMap[projectID] = append(userMap[projectID], userID)
	}
	rows.Close()

	rows, err = db.Query("select id, name, display_name, user_id, description from projects")
	if err != nil {
		return nil, err
	}
	projects := make([]Project, 0)
	for rows.Next() {
		var project Project
		var description sql.NullString
		rows.Scan(&project.ID, &project.Name, &project.DisplayName, &project.UserID, &description)
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
	rows.Close()
	return projects, nil
}

func GetProject(ID string) (Project, error) {
	rows, err := db.Query("select id, name, display_name, user_id, description from projects where id = ?", ID)
	if err != nil {
		return Project{}, err
	}
	if !rows.Next() {
		rows.Close()
		return Project{}, errors.New("project not found")
	}
	var project Project
	var description sql.NullString
	rows.Scan(&project.ID, &project.Name, &project.DisplayName, &project.UserID, &description)
	project.Description = ""
	if description.Valid {
		project.Description = description.String
	}
	rows.Close()

	rows, err = db.Query("select user_id from member where project_id = ?", project.ID)
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

func GetProjectPosts(projectID string, offset, limit int) ([]Post, error) {
	rows, err := db.Query("select posts.id, title, posts.content, thumb_src, posts.user_id, posts.created_at, updated_at, project_id, views, count(comments.id) from (select * from posts where project_id = ? and is_deleted = false order by created_at desc limit ?, ?) posts left join comments on comments.post_id = posts.id group by posts.id", projectID, offset, limit)
	if err != nil {
		return make([]Post, 0), err
	}
	posts := make([]Post, 0)
	for rows.Next() {
		var post Post
		var thumbSrc sql.NullString
		rows.Scan(&post.ID, &post.Title, &post.Content, &thumbSrc, &post.UserID, &post.CreatedAt, &post.UpdatedAt, &post.ProjectID, &post.Views, &post.CommentNum)
		post.ThumbSrc = ""
		if thumbSrc.Valid {
			post.ThumbSrc = thumbSrc.String
		}
		posts = append(posts, post)
	}
	rows.Close()
	return posts, nil
}

func GetPosts(offset, limit int) ([]Post, error) {
	rows, err := db.Query("select posts.id, title, posts.content, thumb_src, posts.user_id, posts.created_at, updated_at, project_id, views, count(comments.id) from (select * from posts where is_deleted = false order by created_at desc limit ?, ?) posts left join comments on comments.post_id = posts.id group by posts.id", offset, limit)
	if err != nil {
		return make([]Post, 0), err
	}
	posts := make([]Post, 0)
	for rows.Next() {
		var post Post
		var thumbSrc sql.NullString
		rows.Scan(&post.ID, &post.Title, &post.Content, &thumbSrc, &post.UserID, &post.CreatedAt, &post.UpdatedAt, &post.ProjectID, &post.Views, &post.CommentNum)
		post.ThumbSrc = ""
		if thumbSrc.Valid {
			post.ThumbSrc = thumbSrc.String
		}
		posts = append(posts, post)
	}
	rows.Close()
	return posts, nil
}

func GetPost(postID string) (Post, error) {
	rows, err := db.Query("select posts.id, title, posts.content, thumb_src, posts.user_id, posts.created_at, updated_at, project_id, views, count(comments.id) from (select * from posts where id = ? and is_deleted = false) posts left join comments on posts.id = comments.post_id group by posts.id", postID)
	if err != nil {
		return Post{}, err
	}
	if !rows.Next() {
		return Post{}, errors.New("post not found")
	}
	var post Post
	var thumbSrc sql.NullString
	rows.Scan(&post.ID, &post.Title, &post.Content, &thumbSrc, &post.UserID, &post.CreatedAt, &post.UpdatedAt, &post.ProjectID, &post.Views, &post.CommentNum)
	post.ThumbSrc = ""
	if thumbSrc.Valid {
		post.ThumbSrc = thumbSrc.String
	}
	return post, nil
}

type Comment struct {
	ID        string `json:"id" form:"id"`
	Content   string `json:"content" form:"content"`
	UserID    string `json:"userId" form:"userId"`
	PostID    string `json:"postId" form:"postId"`
	CreatedAt string `json:"createdAt" form:"createdAt"`
}

func GetPostComments(postID string, offset, limit int) ([]Comment, error) {
	rows, err := db.Query("select id, content, user_id, post_id, created_at from comments where post_id = ? and is_deleted = false order by created_at desc limit ?, ?", postID, offset, limit)
	if err != nil {
		return make([]Comment, 0), err
	}
	comments := make([]Comment, 0)
	for rows.Next() {
		var comment Comment
		rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.CreatedAt)
		comments = append(comments, comment)
	}
	return comments, nil
}

func GetComment(commentID string) (Comment, error) {
	rows, err := db.Query("select id, content, user_id, post_id, created_at from comments where id = ? and is_deleted = false", commentID)
	if err != nil {
		return Comment{}, err
	}
	if !rows.Next() {
		return Comment{}, errors.New("comment not found")
	}
	var comment Comment
	rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.CreatedAt)
	return comment, nil
}

func GetProjectPostById(projectID string, postNumber int) (Post, error) {
	row, err := db.Query("select id from posts where project_id = ? order by created_at asc limit ?, 1", projectID, postNumber)
	if err != nil {
		return Post{}, err
	}
	if !row.Next() {
		row.Close()
		return Post{}, err
	}
	row.Close()

	var postID string
	row.Scan(&postID)
	row, err = db.Query("select posts.id, title, posts.content, thumb_src, posts.user_id, posts.created_at, updated_at, project_id, views, count(comments.id) from (select * from posts where id = ? and is_deleted = false) posts left join comments on posts.id = comments.post_id group by posts.id", postID)

	if err != nil {
		return Post{}, err
	}
	defer row.Close()
	if !row.Next() {
		return Post{}, err
	}

	var post Post
	var thumbSrc sql.NullString
	row.Scan(&post.ID, &post.Title, &post.Content, &thumbSrc, &post.UserID, &post.CreatedAt, &post.UpdatedAt, &post.ProjectID, &post.Views, &post.CommentNum)
	post.ThumbSrc = ""
	if thumbSrc.Valid {
		post.ThumbSrc = thumbSrc.String
	}

	return post, nil
}

func HasAuth(ID string) bool {
	rows, err := db.Query("select auth from users where id = ?", ID)
	if err != nil {
		return false
	}
	defer rows.Close()
	if !rows.Next() {
		return false
	}
	var auth string
	rows.Scan(&auth)
	if auth == "default" {
		go checkProfile(ID)
		return true
	} else {
		return false
	}
}

func HasCommentAuth(ID string) bool {
	rows, err := db.Query("select id from users where id = ?", ID)
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Next()
}

func (post *Post) Insert() error {
	_, err := db.Exec("insert into posts (id, title, content, thumb_src, user_id, created_at, updated_at, project_id, views, is_deleted) value(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", post.ID, post.Title, post.Content, post.ThumbSrc, post.UserID, post.CreatedAt, post.UpdatedAt, post.ProjectID, post.Views, false)
	return err
}

func (project *Project) Insert() error {
	if !regexProjectName.MatchString(project.Name) {
		return errors.New("project name := ^[a-zA-Z0-9_-]+$")
	}
	_, err := db.Exec("insert into projects (id, name, display_name, user_id, description) value(?, ?, ?, ?, ?)", project.ID, project.Name, project.DisplayName, project.UserID, project.Description)
	if err != nil {
		return err
	}
	for _, userID := range project.Member {
		_, err = db.Exec("insert into member (user_id, project_id) value(?, ?)", userID, project.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (comment *Comment) Insert() error {
	_, err := db.Exec("insert into comments (id, content, user_id, post_id, created_at, is_deleted) value(?, ?, ?, ?, ?, ?)", comment.ID, comment.Content, comment.UserID, comment.PostID, comment.CreatedAt, false)
	return err
}

func checkProfile(ID string) {
	rows, err := db.Query("select id from profiles where id = ?", ID)
	if err != nil {
		return
	}
	if !rows.Next() {
		rows.Close()
		db.Query("insert into profiles (id) value(?)", ID)
		return
	}
	rows.Close()
}

func (project *Project) Delete() error {
	_, err := db.Exec("delete from projects where id = ?", project.ID)
	if err != nil {
		return err
	}
	_, err = db.Exec("update posts set is_deleted = true where project_id = ?", project.ID)
	if err != nil {
		return err
	}
	_, err = db.Exec("delete from member where project_id = ?", project.ID)
	if err != nil {
		return err
	}
	_, err = db.Exec("update comments set is_deleted = true where post_id in (select id from posts where project_id = ?)", project.ID)
	return err
}

func (post *Post) Delete() error {
	_, err := db.Exec("update posts set is_deleted = true where id = ?", post.ID)
	if err != nil {
		return err
	}
	_, err = db.Exec("update comments set is_deleted = true where post_id = ?", post.ID)
	return err
}

func (comment *Comment) Delete() error {
	_, err := db.Exec("update comments set is_deleted = true where id = ?", comment.ID)
	return err
}

func (project *Project) Update(invites, removes []string) error {
	if !regexProjectName.MatchString(project.Name) {
		return errors.New("project name := ^[a-zA-Z0-9_-]+$")
	}
	_, err := db.Exec("update projects set name = ?, display_name = ?, user_id = ?, description = ? where id = ?", project.Name, project.DisplayName, project.UserID, project.Description, project.ID)
	for _, userID := range invites {
		db.Exec("insert into member (user_id, project_id) value(?, ?)", userID, project.ID)
	}
	for _, userID := range removes {
		db.Exec("delete from member where user_id = ? and project_id = ?", userID, project.ID)
	}
	return err
}

func (post *Post) Update() error {
	_, err := db.Exec("update posts set title = ?, content = ?, thumb_src = ?, created_at = ?, updated_at = ? where id = ?", post.Title, post.Content, post.ThumbSrc, post.CreatedAt, post.UpdatedAt, post.ID)
	return err
}

func (user *User) Update() error {
	_, err := db.Exec("update profiles set description = ?, twitter_id = ?, github_id = ?, icon_src = ? where id = ?", user.Description, user.TwitterId, user.GithubId, user.IconSrc, user.ID)
	return err
}
