package repositories

import "spyrosmoux/api/internal/models"

type PipelinesRepository interface {
	Save(pipeline *models.Pipeline) (*models.Pipeline, error)
	Update(pipeline *models.Pipeline)
	Delete(pipelineId string) error
	FindById(pipelineId string) (*models.Pipeline, error)
	FindAll() (*[]models.Pipeline, error)
}
