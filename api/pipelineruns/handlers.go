package pipelineruns

import (
	"github.com/gin-gonic/gin"
)

type Handler interface {
	HandleGetPipelineRuns(c *gin.Context)
	HandleUpdatePipelineRun(c *gin.Context)
	HandleUpdatePipelineRunStatus(c *gin.Context)
}
