package project

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupProjectsRouter(db *gorm.DB, router *gin.Engine) {
	projectsRepository := NewProjectsRepositoryImpl(db)
	projectsService := NewProjectsServiceImpl(projectsRepository)
	projectsHandler := NewProjectsHandler(projectsService)

	router.POST("/projects", projectsHandler.AddProject)
	router.GET("/projects", projectsHandler.FindAll)
	router.GET("/projects/:id", projectsHandler.FindProjectById)
	router.DELETE("/projects/:id", projectsHandler.Delete)
}
