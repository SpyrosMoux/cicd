package pipelineruns

import (
	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/cicd/common/dto"
)

type Service interface {
	GetPipelineRuns() dto.ResponseDto
	UpdatePipelineRun(ctx *gin.Context) dto.ResponseDto
	AddPipelineRun(run *PipelineRun) dto.ResponseDto
}
