package pipelines

import "github.com/gin-gonic/gin"

func Routes(route *gin.Engine, pipelineHandler PipelineHandler) {
	routes := route.Group("/app/cicd/api/pipelines")
	{
		routes.POST("", pipelineHandler.HandleCreatePipeline)
		routes.GET("", pipelineHandler.HandleGetAllPipelines)
	}
}
