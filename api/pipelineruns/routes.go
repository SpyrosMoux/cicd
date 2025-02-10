package pipelineruns

import (
	"github.com/gin-gonic/gin"
)

func Routes(rg *gin.RouterGroup, pipelineRunsHandler Handler) {
	pipelineRunsGroup := rg.Group("/runs")
	{
		pipelineRunsGroup.GET("", pipelineRunsHandler.HandleGetPipelineRuns)
		pipelineRunsGroup.PUT("/:id", pipelineRunsHandler.HandleUpdatePipelineRun)
	}
}
