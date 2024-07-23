package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"spyrosmoux/api/internal/models"
	"spyrosmoux/api/internal/services"
)

type RepositoriesHandler struct {
	repositoriesService services.RepositoriesService
}

func NewRepositoriesHandler(repositoriesService services.RepositoriesService) *RepositoriesHandler {
	return &RepositoriesHandler{
		repositoriesService: repositoriesService,
	}
}

func (repositoriesHandler *RepositoriesHandler) CreateRepository(c *gin.Context) {
	var repository models.Repository
	err := c.ShouldBind(&repository)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	createdReporsitory := repositoriesHandler.repositoriesService.Create(repository)
	if createdReporsitory == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"repository": createdReporsitory})
}

func (repositoriesHandler *RepositoriesHandler) UpdateRepository(c *gin.Context) {
	// TODO implement this
	panic("implement me")
}

func (repositoriesHandler *RepositoriesHandler) DeleteRepository(c *gin.Context) {
	id := c.Param("id")
	err := repositoriesHandler.repositoriesService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, "success")
}

func (repositoriesHandler *RepositoriesHandler) FindRepositoryById(c *gin.Context) {
	id := c.Param("id")
	repository, err := repositoriesHandler.repositoriesService.FindById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if repository == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "repository not found"})
	}

	c.JSON(http.StatusOK, gin.H{"repository": repository})
}

func (repositoriesHandler *RepositoriesHandler) FindAllRepositories(c *gin.Context) {
	repositories, err := repositoriesHandler.repositoriesService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"repositories": repositories})
}
