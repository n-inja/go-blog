package main

import (
	"os"

	"github.com/n-inja/go-blog/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	handler.LoadTMPL(r)

	r.GET("blog/", handler.SetTop)

	r.GET("blog/users", handler.SetUsers)
	r.GET("blog/projects", handler.SetProjects)
	r.GET("blog/mypage", handler.SetMyPage)

	r.GET("blog/users/:userID", handler.SetUser)
	r.GET("blog/projects/:projectName", handler.SetProject)

	r.GET("blog/projects/:projectName/:postID/:number", handler.SetPost)
	r.GET("blog/projects/:projectName/:postID", handler.SetPost)

	r.Static("blog/static", os.Getenv("BLOG_STATIC_FILE_PATH")+"/static")

	r.GET("go-blog/api/v1/users", handler.GetAllUsers)
	r.GET("go-blog/api/v1/users/:userID", handler.GetUser)
	r.PUT("go-blog/api/v1/profile", handler.UpdateProfile)

	r.GET("go-blog/api/v1/projects", handler.GetProjects)
	r.GET("go-blog/api/v1/projects/:projectID", handler.GetProject)
	r.POST("go-blog/api/v1/projects", handler.PostProject)
	r.DELETE("go-blog/api/v1/projects/:projectID", handler.DeleteProject)
	r.PUT("go-blog/api/v1/projects/:projectID", handler.UpdateProject)

	r.GET("go-blog/api/v1/users/:userID/posts", handler.GetUserPosts)
	r.GET("go-blog/api/v1/projects/:projectID/posts", handler.GetProjectPosts)
	r.GET("go-blog/api/v1/posts", handler.GetPosts)
	r.GET("go-blog/api/v1/posts/:postID", handler.GetPost)
	r.GET("go-blog/api/v1/projects/:projectID/post", handler.GetProjectPostById)
	r.POST("go-blog/api/v1/projects/:projectID/posts", handler.PostPost)
	r.DELETE("go-blog/api/v1/posts/:postID", handler.DeletePost)
	r.PUT("go-blog/api/v1/posts/:postID", handler.UpdatePost)

	r.GET("go-blog/api/v1/posts/:postID/comments", handler.GetPostComments)
	r.GET("go-blog/api/v1/comments/:commentID", handler.GetComment)
	r.POST("go-blog/api/v1/posts/:postID/comments", handler.PostComment)
	r.DELETE("go-blog/api/v1/comments/:commentID", handler.DeleteComment)

	r.Run(":" + os.Getenv("GO_BLOG_PORT"))
}
