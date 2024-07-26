package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"spyrosmoux/api/internal/handlers"
	"spyrosmoux/api/internal/repositories"
	"spyrosmoux/api/internal/services"
)

func SetupPipelinesRouter(db *gorm.DB, router *gin.Engine) {
	pipelinesRepository := repositories.NewPipelinesRepositoryImpl(db)
	pipelinesService := services.NewPipelinesServiceImpl(pipelinesRepository)
	pipelinesHandler := handlers.NewPipelinesHandler(pipelinesService)

	router.POST("/pipelines", pipelinesHandler.CreatePipeline)
	router.GET("/pipelines", pipelinesHandler.FindAllPipelines)
	router.GET("/pipelines/:id", pipelinesHandler.FindPipelineById)
	router.DELETE("/pipelines/:id", pipelinesHandler.DeletePipeline)
}
