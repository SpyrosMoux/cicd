package pipeline

import (
	"spyrosmoux/api/internal/repository"
)

type PipelinesService interface {
	AutoDiscoverPipelines(repository *repository.Repository) ([]Pipeline, error)
	Create(pipeline *Pipeline) (*Pipeline, error)
	Update(pipeline *Pipeline) (*Pipeline, error)
	Delete(pipelineId string) error
	FindById(pipelineId string) (*Pipeline, error)
	FindAll() (*[]Pipeline, error)
}
