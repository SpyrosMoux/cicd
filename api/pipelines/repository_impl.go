package pipelines

import (
	"gorm.io/gorm"
)

type pipelineRepository struct {
	db *gorm.DB
}

func NewPipelineRepository(db *gorm.DB) PipelineRepository {
	return &pipelineRepository{db: db}
}

func (pipelineRepository pipelineRepository) SavePipeline(pipeline *Pipeline) error {
	result := pipelineRepository.db.Save(&pipeline)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (pipelineRepository pipelineRepository) FindAllPipelines() ([]Pipeline, error) {
	var pipelines []Pipeline
	result := pipelineRepository.db.Find(&pipelines)
	if result.Error != nil {
		return nil, result.Error
	}
	return pipelines, nil
}
