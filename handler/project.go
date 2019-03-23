package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n-inja/go-blog/model"
	"github.com/n-inja/go-blog/utils"
	"github.com/rs/xid"
)

func GetProjects(c *gin.Context) {
	projects, err := model.GetProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func GetProject(c *gin.Context) {
	projectID := c.Param("projectID")
	project, err := model.GetProject(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	c.JSON(http.StatusOK, project)
}

type projectForm struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description"`
}

func PostProject(c *gin.Context) {
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

	project := model.Project{ID: xid.New().String(), Name: projectForm.Name, UserID: ID, Member: []string{ID}, Description: projectForm.Description}
	err = project.Insert()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, project)
}

func DeleteProject(c *gin.Context) {
	projectID := c.Param("projectID")
	ID := c.GetHeader("id")
	project, err := model.GetProject(projectID)
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

type updateProjectForm struct {
	NewName        string   `json:"newName" form:"newName"`
	NewDisplayName string   `json:"newDisplayName" form:"newDisplayName"`
	NewUserID      string   `json:"newUserId" form:"newUserId"`
	NewDescription string   `json:"newDescription" form:"newDescription"`
	Invites        []string `json:"invites" form:"invites"`
	Removes        []string `json:"removes" form:"removes"`
}

func UpdateProject(c *gin.Context) {
	projectID := c.Param("projectID")
	ID := c.GetHeader("id")
	project, err := model.GetProject(projectID)
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
	if body.NewDisplayName != "" {
		project.DisplayName = body.NewDisplayName
	}
	err = project.Update(invites, removes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, project)
}
