package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"./utils"

	"github.com/gin-gonic/gin"
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
