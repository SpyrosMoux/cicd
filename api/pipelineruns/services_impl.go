package pipelineruns

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/cicd/common/dto"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (svc *service) GetPipelineRuns() dto.ResponseDto {
	pipelineRuns, err := svc.repo.FindAll()
	if err != nil {
		return dto.NewResponseDto(http.StatusInternalServerError, "", err.Error(), "")
	}
	return dto.NewResponseDto(http.StatusOK, "", "", pipelineRuns)
}

func (svc *service) UpdatePipelineRun(ctx *gin.Context) dto.ResponseDto {
	runId := ctx.Param("id")

	var run *PipelineRun
	err := ctx.ShouldBindJSON(&run)
	if err != nil {
		return dto.NewResponseDto(http.StatusInternalServerError, "", err.Error(), "")
	}

	savedRun, err := svc.repo.FindById(runId)
	if err != nil {
		return dto.NewResponseDto(http.StatusInternalServerError, "", err.Error(), "")
	}

	savedRun.Status = run.Status

	if run.TimeStarted != 0 {
		savedRun.TimeStarted = run.TimeStarted
	}

	if run.TimeEnded != 0 {
		savedRun.TimeEnded = run.TimeEnded
	}

	updatedRun, err := svc.repo.Update(savedRun)
	if err != nil {
		return dto.NewResponseDto(http.StatusInternalServerError, "", err.Error(), "")
	}

	return dto.NewResponseDto(http.StatusOK, "", "", updatedRun)
}

func (svc *service) AddPipelineRun(run *PipelineRun) dto.ResponseDto {
	err := svc.repo.Create(run)
	if err != nil {
		return dto.NewResponseDto(http.StatusInternalServerError, "", err.Error(), "")
	}
	return dto.NewResponseDto(http.StatusCreated, "Pipeline run created", "", "")
}
