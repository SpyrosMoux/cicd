package services

import (
	"github.com/spyrosmoux/cicd/api/entities"
	"github.com/spyrosmoux/cicd/api/repositories"
)

type PipelineRunsService interface {
	GetPipelineRuns() (*[]entities.PipelineRun, error)
	UpdatePipelineRun(runId string, run *entities.PipelineRun) (*entities.PipelineRun, error)
	UpdatePipelineRunStatus(runId string, status *entities.Status) (*entities.PipelineRun, error)
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

func (svc *pipelineRunsService) UpdatePipelineRun(runId string, run *entities.PipelineRun) (*entities.PipelineRun, error) {
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

func (svc *pipelineRunsService) UpdatePipelineRunStatus(runId string, status *entities.Status) (*entities.PipelineRun, error) {
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
