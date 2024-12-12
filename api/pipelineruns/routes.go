package pipelineruns

import (
	"github.com/gin-gonic/gin"
)

func Routes(route *gin.Engine, pipelineRunsHandler Handler) {
	routes := route.Group("/app/cicd/api/runs")
	{
		routes.GET("", pipelineRunsHandler.HandleGetPipelineRuns)
		routes.PUT("/:id", pipelineRunsHandler.HandleUpdatePipelineRun)
	}
}
