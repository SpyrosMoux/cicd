package services

import "spyrosmoux/api/internal/models"

type PipelinesService interface {
	AutoDiscoverPipelines(repository *models.Repository) ([]models.Pipeline, error)
	Create(pipeline *models.Pipeline) (*models.Pipeline, error)
	Update(pipeline *models.Pipeline) (*models.Pipeline, error)
	Delete(pipelineId string) error
	FindById(pipelineId string) (*models.Pipeline, error)
	FindAll() (*[]models.Pipeline, error)
}
