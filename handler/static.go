package handler

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/n-inja/go-blog/model"
)

var hostname string

func summarize(content string) string {
	length := len(content)
	if length > 32 {
		length = 32
	}
	return strings.Replace(content[0:length], "\"", "", 0)
}

func returnNotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "index.tmpl", gin.H{
		"url":         "https://" + hostname + "/blog/",
		"title":       "CLOG",
		"description": "",
		"imageURL":    "https://" + hostname + "/static/favicon.png",
	})
}

func LoadTMPL(e *gin.Engine) {
	hostname = os.Getenv("HOSTNAME")

	e.LoadHTMLGlob(os.Getenv("BLOG_STATIC_FILE_PATH") + "/*.tmpl")
}

func SetTop(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"url":         "https://" + hostname + "/blog/",
		"title":       "CLOG",
		"description": "",
		"imageURL":    "https://" + hostname + "/static/favicon.png",
	})
}

func SetUsers(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"url":         "https://" + hostname + "/blog/",
		"title":       "CLOG",
		"description": "",
		"imageURL":    "https://" + hostname + "/static/favicon.png",
	})
}

func SetProjects(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"url":         "https://" + hostname + "/blog/projects",
		"title":       "CLOG",
		"description": "",
		"imageURL":    "https://" + hostname + "/static/favicon.png",
	})
}

func SetMyPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"url":         "https://" + hostname + "/blog/mypage",
		"title":       "CLOG",
		"description": "",
		"imageURL":    "https://" + hostname + "/static/favicon.png",
	})
}

func SetUser(c *gin.Context) {
	userID := c.Param("userID")

	user, err := model.GetUser(userID)
	if err != nil {
		returnNotFound(c)
		return
	}
	imageURL := user.IconSrc
	if imageURL != "" {
		imageURL = "https://" + hostname + "/static/favicon.png"
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"url":         "https://" + hostname + "/blog/users/" + user.Name,
		"title":       user.Name,
		"description": summarize(user.Description),
		"imageURL":    imageURL,
	})
}

func SetProject(c *gin.Context) {
	projectName := c.Param("projectName")

	project, err := model.GetProjectByName(projectName)
	if err != nil {
		returnNotFound(c)
		return
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"url":         "https://" + hostname + "/blog/projects/" + projectName,
		"title":       project.Name,
		"description": summarize(project.Description),
		"imageURL":    "https://" + hostname + "/static/favicon.png",
	})
}

func SetPost(c *gin.Context) {
	postID := c.Param("postID")
	projectName := c.Param("projectName")

	if postID != "posts" {
		post, err := model.GetPost(postID)
		if err != nil {
			returnNotFound(c)
			return
		}
		c.Redirect(http.StatusMovedPermanently, "https://"+hostname+"/blog/projects/"+projectName+"/posts/"+strconv.Itoa(post.Number))
		return
	}

	number, err := strconv.Atoi(c.Param("number"))

	if err != nil {
		returnNotFound(c)
		return
	}

	project, err := model.GetProjectByName(projectName)

	if err != nil {
		returnNotFound(c)
		return
	}

	post, err := model.GetProjectPostById(project.ID, number)
	if err != nil {
		returnNotFound(c)
		return
	}

	imageURL := post.ThumbSrc
	if imageURL == "" {
		imageURL = "https://" + hostname + "/static/favicon.png"
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"url":         "https://" + hostname + "/blog/projects/" + projectName + "/post/" + strconv.Itoa(number),
		"title":       project.Name + " - " + post.Title,
		"description": summarize(post.Content),
		"imageURL":    imageURL,
	})
}
