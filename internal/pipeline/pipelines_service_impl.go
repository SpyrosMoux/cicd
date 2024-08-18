package pipeline

import (
	"spyrosmoux/api/internal/repository"
)

type PipelinesServiceImpl struct {
	PipelinesRepository PipelinesRepository
}

func NewPipelinesServiceImpl(pipelinesRepository PipelinesRepository) PipelinesService {
	return &PipelinesServiceImpl{PipelinesRepository: pipelinesRepository}
}

func (p PipelinesServiceImpl) AutoDiscoverPipelines(repository *repository.Repository) ([]Pipeline, error) {
	// TODO implement me
	panic("implement me")
}

func (p PipelinesServiceImpl) Create(pipeline *Pipeline) (*Pipeline, error) {
	newPipeline, err := p.PipelinesRepository.Save(pipeline)
	if err != nil {
		return nil, err
	}

	return newPipeline, nil
}

func (p PipelinesServiceImpl) Update(pipeline *Pipeline) (*Pipeline, error) {
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

func (p PipelinesServiceImpl) FindById(pipelineId string) (*Pipeline, error) {
	pipeline, err := p.PipelinesRepository.FindById(pipelineId)
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

func (p PipelinesServiceImpl) FindAll() (*[]Pipeline, error) {
	pipelines, err := p.PipelinesRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return pipelines, nil
}
