package main

import (
	"fmt"
	"os"

	"./utils"
	"github.com/n-inja/blog/router"

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

	r := gin.Default()
	r.GET("go-blog/api/v1/users", router.GetAllUsers)
	r.GET("go-blog/api/v1/users/:userID", router.GetUser)
	r.PUT("go-blog/api/v1/profile", router.UpdateProfile)

	r.GET("go-blog/api/v1/projects", router.GetProjects)
	r.GET("go-blog/api/v1/projects/:projectID", router.GetProject)
	r.POST("go-blog/api/v1/projects", router.PostProject)
	r.DELETE("go-blog/api/v1/projects/:projectID", router.DeleteProject)
	r.PUT("go-blog/api/v1/projects/:projectID", router.UpdateProject)

	r.GET("go-blog/api/v1/users/:userID/posts", router.GetUserPosts)
	r.GET("go-blog/api/v1/projects/:projectID/posts", router.GetProjectPosts)
	r.GET("go-blog/api/v1/posts", router.GetPosts)
	r.GET("go-blog/api/v1/posts/:postID", router.GetPost)
	r.GET("go-blog/api/v1/projects/:projectID/post", router.GetProjectPostById)
	r.POST("go-blog/api/v1/projects/:projectID/posts", router.PostPost)
	r.DELETE("go-blog/api/v1/posts/:postID", router.DeletePost)
	r.PUT("go-blog/api/v1/posts/:postID", router.UpdatePost)

	r.GET("go-blog/api/v1/posts/:postID/comments", router.GetPostComments)
	r.GET("go-blog/api/v1/comments/:commentID", router.GetComment)
	r.POST("go-blog/api/v1/posts/:postID/comments", router.PostComment)
	r.DELETE("go-blog/api/v1/comments/:commentID", router.DeleteComment)

	r.Run(":" + os.Getenv("GO_BLOG_PORT"))
}
