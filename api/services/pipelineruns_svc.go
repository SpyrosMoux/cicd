package services

import (
	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/cicd/api/entities"
	"github.com/spyrosmoux/cicd/api/repositories"
)

type PipelineRunsService interface {
	GetPipelineRuns() (*[]entities.PipelineRun, error)
	UpdatePipelineRun(ctx *gin.Context) (*entities.PipelineRun, error)
	UpdatePipelineRunStatus(ctx *gin.Context) (*entities.PipelineRun, error)
	AddPipelineRun(run *entities.PipelineRun) error
}

type pipelineRunsService struct {
	repo repositories.PipelineRunsRepository
}

func NewPipelineRunsService(repo repositories.PipelineRunsRepository) PipelineRunsService {
	return &pipelineRunsService{repo: repo}
}

func (svc *pipelineRunsService) GetPipelineRuns() (*[]entities.PipelineRun, error) {
	pipelineRuns, err := svc.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return pipelineRuns, nil
}

func (svc *pipelineRunsService) UpdatePipelineRun(ctx *gin.Context) (*entities.PipelineRun, error) {
	runId := ctx.Param("id")

	var run *entities.PipelineRun
	err := ctx.ShouldBindJSON(&run)
	if err != nil {
		return nil, err
	}

	savedRun, err := svc.repo.FindById(runId)
	if err != nil {
		return nil, err
	}

	savedRun.Status = run.Status
	savedRun.Repository = run.Repository
	savedRun.Branch = run.Branch
	savedRun.TimeTriggered = run.TimeTriggered
	savedRun.TimeStarted = run.TimeStarted
	savedRun.TimeEnded = run.TimeEnded

	updatedRun, err := svc.repo.Update(savedRun)
	if err != nil {
		return nil, err
	}

	return updatedRun, nil
}

func (svc *pipelineRunsService) UpdatePipelineRunStatus(ctx *gin.Context) (*entities.PipelineRun, error) {
	runId := ctx.Param("id")

	var statusStr entities.StatusDto
	err := ctx.ShouldBindJSON(&statusStr)
	if err != nil {
		return nil, err
	}

	status, err := entities.ParseStatus(statusStr.Status)
	if err != nil {
		return nil, err
	}

	savedRun, err := svc.repo.FindById(runId)
	if err != nil {
		return nil, err
	}

	savedRun.Status = status.String()

	updatedRun, err := svc.repo.Update(savedRun)
	if err != nil {
		return nil, err
	}

	return updatedRun, nil
}

func (svc *pipelineRunsService) AddPipelineRun(run *entities.PipelineRun) error {
	return svc.repo.Create(run)
}
