package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/n-inja/blog/model"
	"github.com/n-inja/blog/utils"
	"github.com/rs/xid"
)

func GetUserPosts(c *gin.Context) {
	userID := c.Param("userID")
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
	posts, err := model.GetUserPosts(userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func GetProjectPosts(c *gin.Context) {
	projectID := c.Param("projectID")
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
	posts, err := model.GetProjectPosts(projectID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func GetPosts(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "3"))
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
	posts, err := model.GetPosts(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func GetPost(c *gin.Context) {
	postID := c.Param("postID")
	post, err := model.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, post)
}

func GetProjectPostById(c *gin.Context) {
	projectID := c.Param("projectID")
	postNumber, err := strconv.Atoi(c.DefaultQuery("id", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id should be number"})
		return
	}

	post, err := model.GetProjectPostById(projectID, postNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, post)
}

type postForm struct {
	Title    string `json:"title" form:"title" binding:"required"`
	Content  string `json:"content" form:"content" binding:"required"`
	ThumbSrc string `json:"thumbSrc" form:"thumbSrc"`
}

func PostPost(c *gin.Context) {
	projectID := c.Param("projectID")
	ID := c.GetHeader("id")
	if !utils.HasAuth(ID) {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	var postForm postForm
	err := c.BindJSON(&postForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	date := time.Now()
	post := model.Post{ID: xid.New().String(), Title: postForm.Title, Content: postForm.Content, ThumbSrc: postForm.ThumbSrc, UserID: ID, CreatedAt: date.Format("2006-01-02 15:04:05"), UpdatedAt: date.Format("2006-01-02 15:04:05"), ProjectID: projectID, Views: 0}
	err = post.Insert()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	postID := c.Param("postID")
	ID := c.GetHeader("id")
	post, err := model.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	project, err := model.GetProject(post.ProjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	if ID != project.UserID && ID != post.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	err = post.Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

type updatePostForm struct {
	NewTitle    string `json:"newTitle" form:"newTitle"`
	NewContent  string `json:"newContent" form:"newContent"`
	NewThumbSrc string `json:"newThumbSrc" form:"newThumbSrd"`
}

func UpdatePost(c *gin.Context) {
	postID := c.Param("postID")
	ID := c.GetHeader("id")
	post, err := model.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	if ID != post.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	var body updatePostForm
	err = c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	if body.NewContent != "" {
		post.Content = body.NewContent
	}
	if body.NewTitle != "" {
		post.Title = body.NewTitle
	}
	if body.NewThumbSrc != "" {
		post.ThumbSrc = body.NewThumbSrc
	}
	err = post.Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	post.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	c.JSON(http.StatusOK, post)
}
