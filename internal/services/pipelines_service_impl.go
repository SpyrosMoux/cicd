package services

import (
	"spyrosmoux/api/internal/models"
	"spyrosmoux/api/internal/repositories"
)

type PipelinesServiceImpl struct {
	PipelinesRepository repositories.PipelinesRepository
}

func NewPipelinesServiceImpl(pipelinesRepository repositories.PipelinesRepository) PipelinesService {
	return &PipelinesServiceImpl{PipelinesRepository: pipelinesRepository}
}

func (p PipelinesServiceImpl) AutoDiscoverPipelines(repository *models.Repository) ([]models.Pipeline, error) {
	// TODO implement me
	panic("implement me")
}

func (p PipelinesServiceImpl) Create(pipeline *models.Pipeline) (*models.Pipeline, error) {
	newPipeline, err := p.PipelinesRepository.Save(pipeline)
	if err != nil {
		return nil, err
	}

	return newPipeline, nil
}

func (p PipelinesServiceImpl) Update(pipeline *models.Pipeline) (*models.Pipeline, error) {
	//TODO implement me
	panic("implement me")
}

func (p PipelinesServiceImpl) Delete(pipelineId string) error {
	err := p.PipelinesRepository.Delete(pipelineId)
	if err != nil {
		return err
	}

	return nil
}

func (p PipelinesServiceImpl) FindById(pipelineId string) (*models.Pipeline, error) {
	pipeline, err := p.PipelinesRepository.FindById(pipelineId)
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

func (p PipelinesServiceImpl) FindAll() (*[]models.Pipeline, error) {
	pipelines, err := p.PipelinesRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return pipelines, nil
}
