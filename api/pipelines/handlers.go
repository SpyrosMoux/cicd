package pipelines

import "github.com/gin-gonic/gin"

type PipelineHandler interface {
	HandleCreatePipeline(ctx *gin.Context)
	HandleGetAllPipelines(ctx *gin.Context)
	HandleDeletePipeline(ctx *gin.Context)
}
