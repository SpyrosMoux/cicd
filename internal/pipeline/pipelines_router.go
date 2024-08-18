package pipeline

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupPipelinesRouter(db *gorm.DB, router *gin.Engine) {
	pipelinesRepository := NewPipelinesRepositoryImpl(db)
	pipelinesService := NewPipelinesServiceImpl(pipelinesRepository)
	pipelinesHandler := NewPipelinesHandler(pipelinesService)

	router.POST("/pipelines", pipelinesHandler.CreatePipeline)
	router.GET("/pipelines", pipelinesHandler.FindAllPipelines)
	router.GET("/pipelines/:id", pipelinesHandler.FindPipelineById)
	router.DELETE("/pipelines/:id", pipelinesHandler.DeletePipeline)
}
