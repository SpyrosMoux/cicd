package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/cicd/api/handlers"
)

func PipelineRuns(route *gin.Engine, pipelineRunsHandler handlers.PipelineRunsHandler) {
	routes := route.Group("/app/cicd/api/runs")
	{
		routes.GET("", pipelineRunsHandler.HandleGetPipelineRuns)
		routes.POST("/:id", pipelineRunsHandler.HandleUpdatePipelineRun) // TODO(@SpyrosMoux) this should be PUT
		routes.PUT("/:id", pipelineRunsHandler.HandleUpdatePipelineRunStatus)
	}
}
