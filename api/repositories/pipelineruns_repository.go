package repositories

import (
	"github.com/spyrosmoux/cicd/api/config"
	"github.com/spyrosmoux/cicd/api/entities"
	"gorm.io/gorm"
)

type PipelineRunsRepository interface {
	FindById(pipelineRunId string) (*entities.PipelineRun, error)
	FindAll() (*[]entities.PipelineRun, error)
	Update(run *entities.PipelineRun) (*entities.PipelineRun, error)
	Create(run *entities.PipelineRun) error
}

type pipelineRunsRepository struct {
	db *gorm.DB
}

func NewPipelineRunsRepository(db *gorm.DB) PipelineRunsRepository {
	return &pipelineRunsRepository{db: db}
}

func (repo *pipelineRunsRepository) FindById(id string) (*entities.PipelineRun, error) {
	var run entities.PipelineRun
	result := config.DB.Where("id = ?", id).First(&run)
	if result.Error != nil {
		return &entities.PipelineRun{}, result.Error
	}
	return &run, nil
}

func (repo *pipelineRunsRepository) FindAll() (*[]entities.PipelineRun, error) {
	var runs *[]entities.PipelineRun
	result := config.DB.Find(&runs)
	if result.Error != nil {
		return &[]entities.PipelineRun{}, result.Error
	}

	return runs, nil
}

func (repo *pipelineRunsRepository) Update(run *entities.PipelineRun) (*entities.PipelineRun, error) {
	result := config.DB.Save(&run)
	if result.Error != nil {
		return nil, result.Error
	}
	return run, nil
}

func (repo *pipelineRunsRepository) Create(run *entities.PipelineRun) error {
	result := config.DB.Create(run)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
