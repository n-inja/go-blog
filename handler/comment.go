package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/n-inja/go-blog/model"
	"github.com/n-inja/go-blog/utils"
	"github.com/rs/xid"
)

func GetPostComments(c *gin.Context) {
	postID := c.Param("postID")
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}
	if limit < 0 || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	comments, err := model.GetPostComments(postID, offset, limit)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func GetComment(c *gin.Context) {
	commentID := c.Param("commentID")
	comment, err := model.GetComment(commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, comment)
}

type commentForm struct {
	Content string `json:"content" form:"content" binding:"required"`
}

func PostComment(c *gin.Context) {
	postID := c.Param("postID")
	ID := c.GetHeader("id")
	if !utils.HasCommentAuth(ID) {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	var commentForm commentForm
	err := c.BindJSON(&commentForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	comment := model.Comment{ID: xid.New().String(), Content: commentForm.Content, UserID: ID, PostID: postID, CreatedAt: time.Now().Format("2006-01-02 15:04:05")}
	err = comment.Insert()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, comment)
}

func DeleteComment(c *gin.Context) {
	commentID := c.Param("commentID")
	ID := c.GetHeader("id")
	comment, err := model.GetComment(commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	if ID != comment.ID {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	err = comment.Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
