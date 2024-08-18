package project

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProjectsHandler struct {
	projectsService ProjectsService
}

func NewProjectsHandler(projectsService ProjectsService) *ProjectsHandler {
	return &ProjectsHandler{
		projectsService: projectsService,
	}
}

func (projectsHandler *ProjectsHandler) AddProject(c *gin.Context) {
	var project Project
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

	if project == nil {
		c.IndentedJSON(http.StatusNotFound, nil)
		return
	}

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
