package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"spyrosmoux/api/internal/handlers"
	"spyrosmoux/api/internal/repositories"
	"spyrosmoux/api/internal/services"
)

func SetupRepositoriesRouter(db *gorm.DB, router *gin.Engine) {
	repositoriesRepository := repositories.NewRepositoriesRepositoryImpl(db)
	repositoriesService := services.NewRepositoriesServiceImpl(repositoriesRepository)
	repositoriesHandler := handlers.NewRepositoriesHandler(repositoriesService)

	router.POST("/repositories", repositoriesHandler.CreateRepository)
	router.GET("/repositories", repositoriesHandler.FindAllRepositories)
	router.GET("/repositories/:id", repositoriesHandler.FindRepositoryById)
	router.DELETE("/repositories/:id", repositoriesHandler.DeleteRepository)

}
