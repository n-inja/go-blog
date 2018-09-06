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
	router.GET("go-blog/api/v1/users/:userID/posts", getUserPosts)
	router.GET("go-blog/api/v1/projects", getProjects)
	router.GET("go-blog/api/v1/projects/:projectID", getProject)
	router.GET("go-blog/api/v1/projects/:projectID/posts", getProjectPosts)
	router.GET("go-blog/api/v1/posts", getPosts)
	router.GET("go-blog/api/v1/posts/:postID", getPost)
	router.GET("go-blog/api/v1/posts/:postID/comments", getPostComments)
	router.GET("go-blog/api/v1/comments/:commentID", getComment)

	router.POST("go-blog/api/v1/projects/:projectID/posts", postPost)
	router.POST("go-blog/api/v1/projects", postProject)
	router.POST("go-blog/api/v1/posts/:postID/comments", postComment)

	router.DELETE("go-blog/api/v1/projects/:projectID", deleteProject)
	router.DELETE("go-blog/api/v1/posts/:postID", deletePost)
	router.DELETE("go-blog/api/v1/comments/:commentID", deleteComment)

	router.PUT("go-blog/api/v1/projects/:projectID", updateProject)
	router.PUT("go-blog/api/v1/posts/:postID", updatePost)
	router.PUT("go-blog/api/v1/profile", updateProfile)

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

func getUserPosts(c *gin.Context) {
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
	posts, err := utils.GetUserPosts(userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, posts)
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

func getPosts(c *gin.Context) {
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
	posts, err := utils.GetPosts(offset, limit)
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
	err := c.BindJSON(&projectForm)
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

func deleteProject(c *gin.Context) {
	projectID := c.Param("projectID")
	ID := c.GetHeader("id")
	project, err := utils.GetProject(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	if project.UserID != ID {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	err = project.Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func deletePost(c *gin.Context) {
	postID := c.Param("postID")
	ID := c.GetHeader("id")
	post, err := utils.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	project, err := utils.GetProject(post.ProjectID)
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

func deleteComment(c *gin.Context) {
	commentID := c.Param("commentID")
	ID := c.GetHeader("id")
	comment, err := utils.GetComment(commentID)
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

type updateProjectForm struct {
	NewName        string   `json:"newName" form:"newName"`
	NewUserID      string   `json:"newUserId" form:"newUserId"`
	NewDescription string   `json:"newDescription" form:"newDescription"`
	Invites        []string `json:"invites" form:"invites"`
	Removes        []string `json:"removes" form:"removes"`
}

func updateProject(c *gin.Context) {
	projectID := c.Param("projectID")
	ID := c.GetHeader("id")
	project, err := utils.GetProject(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	if project.UserID != ID {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	var body updateProjectForm
	err = c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	memberMap := make(map[string]bool, 0)
	for _, userID := range project.Member {
		memberMap[userID] = true
	}
	invites := make([]string, 0)
	for _, userID := range body.Invites {
		if !memberMap[userID] {
			invites = append(invites, userID)
		}
	}
	removes := make([]string, 0)
	for _, userID := range body.Invites {
		if memberMap[userID] && userID != project.UserID && userID != body.NewUserID {
			removes = append(removes, userID)
		}
	}

	if body.NewUserID != "" {
		belong := false
		for _, userID := range project.Member {
			if userID == body.NewUserID {
				belong = true
			}
		}
		if !belong {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "project owner should belong project",
			})
			return
		}
		project.UserID = body.NewUserID
	}
	if body.NewName != "" {
		project.Name = body.NewName
	}
	if body.NewDescription != "" {
		project.Description = body.NewDescription
	}
	err = project.Update(invites, removes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, project)
}

type UpdatePostForm struct {
	NewTitle    string `json:"newTitle" form:"newTitle"`
	NewContent  string `json:"newContent" form:"newContent"`
	NewThumbSrc string `json:"newThumbSrc" form:"newThumbSrd"`
}

func updatePost(c *gin.Context) {
	postID := c.Param("postID")
	ID := c.GetHeader("id")
	post, err := utils.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	if ID != post.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	var body UpdatePostForm
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

type UpdateProfileForm struct {
	NewDescription string `json:"newDescription" form:"newDescription"`
	NewTwitterID   string `json:"newTwitterId" form:"newTwitterId"`
	NewGithubID    string `json:"newGithubId" form:"newGithubId"`
	NewIconSrc     string `json:"newIconSrc" form:"newIconSrc"`
}

func updateProfile(c *gin.Context) {
	ID := c.GetHeader("id")
	user, err := utils.GetUser(ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	var body UpdateProfileForm
	err = c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	if body.NewDescription != "" {
		user.Description = body.NewDescription
	}
	if body.NewGithubID != "" {
		user.GithubId = body.NewGithubID
	}
	if body.NewTwitterID != "" {
		user.TwitterId = body.NewTwitterID
	}
	if body.NewIconSrc != "" {
		user.IconSrc = body.NewIconSrc
	}
	err = user.Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, user)
}
