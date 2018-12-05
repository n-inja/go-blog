package model

import (
	"errors"

	"github.com/n-inja/go-blog/utils"
)

func (comment *Comment) Insert() error {
	utils.Open()

	_, err := utils.DB.Exec("insert into comments (id, content, user_id, post_id, created_at, is_deleted) value(?, ?, ?, ?, ?, ?)", comment.ID, comment.Content, comment.UserID, comment.PostID, comment.CreatedAt, false)
	return err
}

func (comment *Comment) Delete() error {
	utils.Open()

	_, err := utils.DB.Exec("update comments set is_deleted = true where id = ?", comment.ID)
	return err
}

func (comment *Comment) Update() error {
	utils.Open()

	_, err := utils.DB.Exec("update comments set content = ? where id = ?", comment.Content, comment.ID)
	return err
}

type Comment struct {
	ID        string `json:"id" form:"id"`
	Content   string `json:"content" form:"content"`
	UserID    string `json:"userId" form:"userId"`
	PostID    string `json:"postId" form:"postId"`
	CreatedAt string `json:"createdAt" form:"createdAt"`
}

func GetPostComments(postID string, offset, limit int) ([]Comment, error) {
	utils.Open()

	rows, err := utils.DB.Query("select id, content, user_id, post_id, created_at from comments where post_id = ? and is_deleted = false order by created_at desc limit ?, ?", postID, offset, limit)
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
	utils.Open()

	rows, err := utils.DB.Query("select id, content, user_id, post_id, created_at from comments where id = ? and is_deleted = false", commentID)
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
