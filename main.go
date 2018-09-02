package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"./utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func main() {
	// connect database
	databaseAddress := ""
	if os.Getenv("DATABASE_ADDRESS") != "" {
		databaseAddress = os.Getenv("DATABASE_ADDRESS")
	}
	err := utils.Open(os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD"), databaseAddress, os.Getenv("DATABASE_NAME"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer utils.Close()

	router := gin.Default()
	router.GET("go-blog/api/v1/users", getAllUsers)
	router.GET("go-blog/api/v1/users/:userID", getUser)
	router.GET("go-blog/api/v1/projects", getProjects)
	router.GET("go-blog/api/v1/projects/:projectID", getProject)
	router.GET("go-blog/api/v1/projects/:projectID/posts", getProjectPosts)
	router.GET("go-blog/api/v1/posts/:postID", getPost)
	router.GET("go-blog/api/v1/posts/:postID/comments", getPostComments)
	router.GET("go-blog/api/v1/comments/:commentID", getComment)

	router.POST("go-blog/api/v1/projects/:projectID/posts", postPost)
	router.POST("go-blog/api/v1/projects", postProject)
	router.POST("go-blog/api/v1/posts/:postID/comments", postComment)

	router.Run(":" + os.Getenv("GO_BLOG_PORT"))
}

func getAllUsers(c *gin.Context) {
	users, err := utils.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	userID := c.Param("userID")
	user, err := utils.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, user)
}

func getProjects(c *gin.Context) {
	projects, err := utils.GetProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func getProject(c *gin.Context) {
	projectID := c.Param("projectID")
	project, err := utils.GetProject(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, project)
}

func getProjectPosts(c *gin.Context) {
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
	posts, err := utils.GetProjectPosts(projectID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func getPost(c *gin.Context) {
	postID := c.Param("postID")
	post, err := utils.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, post)
}

func getPostComments(c *gin.Context) {
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

	comments, err := utils.GetPostComments(postID, offset, limit)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func getComment(c *gin.Context) {
	commentID := c.Param("commentID")
	comment, err := utils.GetComment(commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, comment)
}

type PostForm struct {
	Title    string `json:"title" form:"title" binding:"required"`
	Content  string `json:"content" form:"content" binding:"required"`
	ThumbSrc string `json:"thumbSrc" form:"thumbSrc"`
}

func postPost(c *gin.Context) {
	projectID := c.Param("projectID")
	ID := c.GetHeader("id")
	if !utils.HasAuth(ID) {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	var postForm PostForm
	err := c.BindJSON(&postForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	date := time.Now()
	post := utils.Post{ID: xid.New().String(), Title: postForm.Title, Content: postForm.Content, ThumbSrc: postForm.ThumbSrc, UserID: ID, CreatedAt: date.Format("2006-01-02 15:04:05"), UpdatedAt: date.Format("2006-01-02 15:04:05"), ProjectID: projectID, Views: 0}
	err = post.Insert()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, post)
}

type projectForm struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description"`
}

func postProject(c *gin.Context) {
	ID := c.GetHeader("id")
	if !utils.HasAuth(ID) {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	var projectForm projectForm
	err := c.BindJSON(projectForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	project := utils.Project{ID: xid.New().String(), Name: projectForm.Name, UserID: ID, Member: []string{ID}, Description: projectForm.Description}
	err = project.Insert()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, project)
}

type CommentForm struct {
	Content string `json:"content" form:"content" binding:"required"`
}

func postComment(c *gin.Context) {
	postID := c.Param("postID")
	ID := c.GetHeader("id")
	if !utils.HasCommentAuth(ID) {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	var commentForm CommentForm
	err := c.BindJSON(&commentForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	comment := utils.Comment{ID: xid.New().String(), Content: commentForm.Content, UserID: ID, PostID: postID, CreatedAt: time.Now().Format("2006-01-02 15:04:05")}
	err = comment.Insert()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, comment)
}
