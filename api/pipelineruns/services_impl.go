package pipelineruns

import (
	"github.com/gin-gonic/gin"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (svc *service) GetPipelineRuns() (*[]PipelineRun, error) {
	pipelineRuns, err := svc.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return pipelineRuns, nil
}

func (svc *service) UpdatePipelineRun(ctx *gin.Context) (*PipelineRun, error) {
	runId := ctx.Param("id")

	var run *PipelineRun
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

func (svc *service) UpdatePipelineRunStatus(ctx *gin.Context) (*PipelineRun, error) {
	runId := ctx.Param("id")

	var statusStr StatusDto
	err := ctx.ShouldBindJSON(&statusStr)
	if err != nil {
		return nil, err
	}

	status, err := ParseStatus(statusStr.Status)
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

func (svc *service) AddPipelineRun(run *PipelineRun) error {
	return svc.repo.Create(run)
}
