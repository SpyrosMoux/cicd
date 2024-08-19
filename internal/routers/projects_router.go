package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"spyrosmoux/api/internal/handlers"
	"spyrosmoux/api/internal/repositories"
	"spyrosmoux/api/internal/services"
)

func SetupProjectsRouter(db *gorm.DB, router *gin.Engine) {
	projectsRepository := repositories.NewProjectsRepositoryImpl(db)
	projectsService := services.NewProjectsServiceImpl(projectsRepository)
	projectsHandler := handlers.NewProjectsHandler(projectsService)

	router.POST("/projects", projectsHandler.AddProject)
	router.GET("/projects", projectsHandler.FindAll)
	router.GET("/projects/:id", projectsHandler.FindProjectById)
	router.DELETE("/projects/:id", projectsHandler.Delete)
}
