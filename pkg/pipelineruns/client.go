package pipelineruns

import (
	"log"

	"github.com/spyrosmoux/api/internal/db"
)

type Client interface {
	GetPipelineRuns() (*[]PipelineRun, error)
	AddPipelineRun(pipelineRun *PipelineRun) error
	UpdatePipelineRun(pipelineRunId string, pipelineRun *PipelineRun) (*PipelineRun, error)
	UpdatePipelineRunStatus(pipelineRunId string, status Status) (*PipelineRun, error)
}

type ClientImpl struct{}

func NewClient() (Client, error) {
	return &ClientImpl{}, nil
}

func (client *ClientImpl) GetPipelineRuns() (*[]PipelineRun, error) {
	var pipelineRuns *[]PipelineRun

	result := db.DB.Find(&pipelineRuns)
	if result.Error != nil {
		return &[]PipelineRun{}, result.Error
	}

	return pipelineRuns, nil
}

func (client *ClientImpl) AddPipelineRun(pipelineRun *PipelineRun) error {
	result := db.DB.Create(pipelineRun)
	if result.Error != nil {
		log.Printf("Error adding pipeline run: %v", result.Error)
		return result.Error
	}
	return nil
}

func (client *ClientImpl) UpdatePipelineRun(pipelineRunId string, pipelineRun *PipelineRun) (*PipelineRun, error) {
	var savedPipelineRun PipelineRun
	result := db.DB.Find(&savedPipelineRun, pipelineRunId)
	if result.Error != nil {
		log.Printf("Error finding pipeline run with id: %v and error: %v", pipelineRunId, result.Error)
		return &PipelineRun{}, result.Error
	}

	savedPipelineRun.Status = pipelineRun.Status
	savedPipelineRun.Repository = pipelineRun.Repository
	savedPipelineRun.Branch = pipelineRun.Branch
	savedPipelineRun.TimeTriggered = pipelineRun.TimeTriggered
	savedPipelineRun.TimeStarted = pipelineRun.TimeStarted
	savedPipelineRun.TimeEnded = pipelineRun.TimeEnded

	result = db.DB.Save(&savedPipelineRun)
	if result.Error != nil {
		log.Printf("Error updating pipeline run: %v", result.Error)
		return &PipelineRun{}, result.Error
	}

	return &savedPipelineRun, nil
}

func (client *ClientImpl) UpdatePipelineRunStatus(pipelineRunId string, status Status) (*PipelineRun, error) {
	var savedPipelineRun PipelineRun
	result := db.DB.Find(&savedPipelineRun, pipelineRunId)
	if result.Error != nil {
		log.Printf("Error finding pipeline run with id: %v and error: %v", pipelineRunId, result.Error)
		return &PipelineRun{}, result.Error
	}

	savedPipelineRun.Status = status.String()

	result = db.DB.Save(savedPipelineRun)
	if result.Error != nil {
		log.Printf("Error updating pipeline run: %v", result.Error)
		return &PipelineRun{}, result.Error
	}

	return &savedPipelineRun, nil
}
