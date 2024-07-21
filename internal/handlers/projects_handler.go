package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"spyrosmoux/api/internal/models"
	"spyrosmoux/api/internal/services"
)

type ProjectsHandler struct {
	projectsService services.ProjectsService
}

func NewProjectsHandler(projectsService services.ProjectsService) *ProjectsHandler {
	return &ProjectsHandler{
		projectsService: projectsService,
	}
}

func (projectsHandler *ProjectsHandler) AddProject(c *gin.Context) {
	var project models.Project
	err := c.BindJSON(&project)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	createdProject := projectsHandler.projectsService.Create(project)
	if createdProject == nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.IndentedJSON(http.StatusCreated, createdProject)
}

func (projectsHandler *ProjectsHandler) FindProjectById(c *gin.Context) {
	id := c.Param("id")
	project := projectsHandler.projectsService.FindById(id)

	c.IndentedJSON(http.StatusOK, project)
}

func (projectsHandler *ProjectsHandler) FindAll(c *gin.Context) {
	projects := projectsHandler.projectsService.FindAll()

	c.IndentedJSON(http.StatusOK, projects)
}

func (projectsHandler *ProjectsHandler) Delete(c *gin.Context) {
	err := projectsHandler.projectsService.Delete(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, "success")
}
