package repositories

import (
	"gorm.io/gorm"
	"spyrosmoux/api/internal/models"
)

type PipelinesRepositoryImpl struct {
	Db *gorm.DB
}

func NewPipelinesRepositoryImpl(Db *gorm.DB) PipelinesRepository {
	return &PipelinesRepositoryImpl{Db: Db}
}

func (p PipelinesRepositoryImpl) Save(pipeline *models.Pipeline) (*models.Pipeline, error) {
	result := p.Db.Create(&pipeline)
	if result.Error != nil {
		return nil, result.Error
	}

	return pipeline, nil
}

func (p PipelinesRepositoryImpl) Update(pipeline *models.Pipeline) {
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

func (p PipelinesRepositoryImpl) FindById(pipelineId string) (*models.Pipeline, error) {
	pipeline := &models.Pipeline{}
	result := p.Db.Find(&pipeline, "id = ?", pipelineId)
	if result.Error != nil {
		return nil, result.Error
	}

	return pipeline, nil
}

func (p PipelinesRepositoryImpl) FindAll() (*[]models.Pipeline, error) {
	pipelines := []models.Pipeline{}
	result := p.Db.Find(&pipelines)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pipelines, nil
}
