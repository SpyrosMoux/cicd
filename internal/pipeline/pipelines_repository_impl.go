package pipeline

import (
	"gorm.io/gorm"
)

type PipelinesRepositoryImpl struct {
	Db *gorm.DB
}

func NewPipelinesRepositoryImpl(Db *gorm.DB) PipelinesRepository {
	return &PipelinesRepositoryImpl{Db: Db}
}

func (p PipelinesRepositoryImpl) Save(pipeline *Pipeline) (*Pipeline, error) {
	result := p.Db.Create(&pipeline)
	if result.Error != nil {
		return nil, result.Error
	}

	return pipeline, nil
}

func (p PipelinesRepositoryImpl) Update(pipeline *Pipeline) {
	//TODO implement me
	panic("implement me")
}

func (p PipelinesRepositoryImpl) Delete(pipelineId string) error {
	pipeline, err := p.FindById(pipelineId)
	if err != nil {
		return err
	}

	result := p.Db.Delete(&pipeline)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (p PipelinesRepositoryImpl) FindById(pipelineId string) (*Pipeline, error) {
	pipeline := &Pipeline{}
	result := p.Db.Find(&pipeline, "id = ?", pipelineId)
	if result.Error != nil {
		return nil, result.Error
	}

	return pipeline, nil
}

func (p PipelinesRepositoryImpl) FindAll() (*[]Pipeline, error) {
	pipelines := []Pipeline{}
	result := p.Db.Find(&pipelines)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pipelines, nil
}
