package model

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/n-inja/go-blog/utils"
)

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
	Number     int    `json:"number" form:"number"`
}

func (post *Post) Insert() error {
	return utils.Transact(func(tx *sql.Tx) error {
		err := tx.QueryRow("select count(*) from posts where project_id = ?", post.ProjectID).Scan(&post.Number)
		if err != nil {
			fmt.Println("a")
			return err
		}
		_, err = tx.Exec("insert into posts (id, title, content, thumb_src, user_id, number, project_id, views, is_deleted) value(?, ?, ?, ?, ?, ?, ?, ?, ?)", post.ID, post.Title, post.Content, post.ThumbSrc, post.UserID, post.Number, post.ProjectID, post.Views, false)
		fmt.Println("b")
		return err
	})
}

func (post *Post) Delete() error {
	return utils.Transact(func(tx *sql.Tx) error {
		_, err := tx.Exec("update posts set is_deleted = true where id = ?", post.ID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("update comments set is_deleted = true where post_id = ?", post.ID)
		return err
	})
}

func (post *Post) Update() error {
	_, err := utils.DB.Exec("update posts set title = ?, content = ?, thumb_src = ? where id = ?", post.Title, post.Content, post.ThumbSrc, post.ID)
	return err
}

func GetUserPosts(userID string, offset, limit int) ([]Post, error) {
	rows, err := utils.DB.Query("select posts.id, title, posts.content, thumb_src, posts.user_id, posts.number, posts.created_at, updated_at, project_id, views, count(comments.id) from (select * from posts where user_id = ? and is_deleted = false order by created_at desc limit ?, ?) posts left join comments on comments.post_id = posts.id group by posts.id", userID, offset, limit)
	if err != nil {
		return make([]Post, 0), err
	}
	posts := make([]Post, 0)
	for rows.Next() {
		var post Post
		var thumbSrc sql.NullString
		rows.Scan(&post.ID, &post.Title, &post.Content, &thumbSrc, &post.UserID, &post.Number, &post.CreatedAt, &post.UpdatedAt, &post.ProjectID, &post.Views, &post.CommentNum)
		post.ThumbSrc = ""
		if thumbSrc.Valid {
			post.ThumbSrc = thumbSrc.String
		}
		posts = append(posts, post)
	}
	rows.Close()
	return posts, nil
}

func GetProjectPosts(projectID string, offset, limit int) ([]Post, error) {
	rows, err := utils.DB.Query("select posts.id, title, posts.content, thumb_src, posts.user_id, posts.number, posts.created_at, updated_at, project_id, views, count(comments.id) from (select * from posts where project_id = ? and is_deleted = false order by created_at desc limit ?, ?) posts left join comments on comments.post_id = posts.id group by posts.id", projectID, offset, limit)
	if err != nil {
		return make([]Post, 0), err
	}
	posts := make([]Post, 0)
	for rows.Next() {
		var post Post
		var thumbSrc sql.NullString
		rows.Scan(&post.ID, &post.Title, &post.Content, &thumbSrc, &post.UserID, &post.Number, &post.CreatedAt, &post.UpdatedAt, &post.ProjectID, &post.Views, &post.CommentNum)
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
	rows, err := utils.DB.Query("select posts.id, title, posts.content, thumb_src, posts.user_id, posts.number, posts.created_at, updated_at, project_id, views, count(comments.id) from (select * from posts where is_deleted = false order by created_at desc limit ?, ?) posts left join comments on comments.post_id = posts.id group by posts.id", offset, limit)
	if err != nil {
		return make([]Post, 0), err
	}
	posts := make([]Post, 0)
	for rows.Next() {
		var post Post
		var thumbSrc sql.NullString
		rows.Scan(&post.ID, &post.Title, &post.Content, &thumbSrc, &post.UserID, &post.Number, &post.CreatedAt, &post.UpdatedAt, &post.ProjectID, &post.Views, &post.CommentNum)
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
	var post Post
	var thumbSrc sql.NullString
	err := utils.DB.QueryRow("select posts.id, title, posts.content, thumb_src, posts.user_id, posts.number, posts.created_at, updated_at, project_id, views, count(comments.id) from (select * from posts where id = ? and is_deleted = false) posts left join comments on posts.id = comments.post_id group by posts.id", postID).Scan(&post.ID, &post.Title, &post.Content, &thumbSrc, &post.UserID, &post.Number, &post.CreatedAt, &post.UpdatedAt, &post.ProjectID, &post.Views, &post.CommentNum)
	if err != nil {
		return Post{}, err
	}
	post.ThumbSrc = ""
	if thumbSrc.Valid {
		post.ThumbSrc = thumbSrc.String
	}

	return post, nil
}

func GetProjectPostById(projectID string, postNumber int) (Post, error) {
	var post Post
	var thumbSrc sql.NullString
	err := utils.DB.QueryRow("select posts.id, title, posts.content, thumb_src, posts.user_id, posts.number, posts.created_at, updated_at, project_id, views, count(comments.id) from (select * from posts where project_id = ? and number = ? and is_deleted = false) posts left join comments on posts.id = comments.post_id group by posts.id", projectID, postNumber).Scan(&post.ID, &post.Title, &post.Content, &thumbSrc, &post.UserID, &post.Number, &post.CreatedAt, &post.UpdatedAt, &post.ProjectID, &post.Views, &post.CommentNum)
	if err != nil {
		return Post{}, errors.New("db error")
	}

	post.ThumbSrc = ""
	if thumbSrc.Valid {
		post.ThumbSrc = thumbSrc.String
	}
	post.Number = postNumber

	return post, nil
}
