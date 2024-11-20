package pipelineruns

import (
	"github.com/gin-gonic/gin"
)

type Service interface {
	GetPipelineRuns() (*[]PipelineRun, error)
	UpdatePipelineRun(ctx *gin.Context) (*PipelineRun, error)
	UpdatePipelineRunStatus(ctx *gin.Context) (*PipelineRun, error)
	AddPipelineRun(run *PipelineRun) error
}
