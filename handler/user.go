package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n-inja/go-blog/model"
)

func GetAllUsers(c *gin.Context) {
	users, err := model.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	userID := c.Param("userID")
	user, err := model.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, user)
}

type updateProfileForm struct {
	NewDescription string `json:"newDescription" form:"newDescription"`
	NewTwitterID   string `json:"newTwitterId" form:"newTwitterId"`
	NewGithubID    string `json:"newGithubId" form:"newGithubId"`
	NewIconSrc     string `json:"newIconSrc" form:"newIconSrc"`
}

func UpdateProfile(c *gin.Context) {
	ID := c.GetHeader("id")
	user, err := model.GetUser(ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	var body updateProfileForm
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
