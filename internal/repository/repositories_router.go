package repository

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRepositoriesRouter(db *gorm.DB, router *gin.Engine) {
	repositoriesRepository := NewRepositoriesRepositoryImpl(db)
	repositoriesService := NewRepositoriesServiceImpl(repositoriesRepository)
	repositoriesHandler := NewRepositoriesHandler(repositoriesService)

	router.POST("/repositories", repositoriesHandler.CreateRepository)
	router.GET("/repositories", repositoriesHandler.FindAllRepositories)
	router.GET("/repositories/:id", repositoriesHandler.FindRepositoryById)
	router.DELETE("/repositories/:id", repositoriesHandler.DeleteRepository)

}
