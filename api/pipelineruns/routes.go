package pipelineruns

import (
	"github.com/gin-gonic/gin"
)

func Routes(route *gin.Engine, pipelineRunsHandler Handler) {
	routes := route.Group("/app/cicd/api/runs")
	{
		routes.GET("", pipelineRunsHandler.HandleGetPipelineRuns)
		routes.POST("/:id", pipelineRunsHandler.HandleUpdatePipelineRun) // TODO(@SpyrosMoux) this should be PUT
		routes.PUT("/:id", pipelineRunsHandler.HandleUpdatePipelineRunStatus)
	}
}
