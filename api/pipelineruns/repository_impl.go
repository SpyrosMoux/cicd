package pipelineruns

import (
	"fmt"

	"github.com/spyrosmoux/cicd/api/config"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (repo *repository) FindById(id string) (*PipelineRun, error) {
	var run PipelineRun
	result := config.DB.Where("id = ?", id).First(&run)
	if result.Error != nil {
		return &PipelineRun{}, result.Error
	}
	return &run, nil
}

func (repo *repository) FindAll() (*[]PipelineRun, error) {
	var runs *[]PipelineRun
	result := config.DB.Find(&runs)
	if result.Error != nil {
		return &[]PipelineRun{}, fmt.Errorf("unable to fetch pipeline runs from db, err=%s", result.Error)
	}

	return runs, nil
}

func (repo *repository) Update(run *PipelineRun) (*PipelineRun, error) {
	result := config.DB.Save(&run)
	if result.Error != nil {
		return nil, result.Error
	}
	return run, nil
}

func (repo *repository) Create(run *PipelineRun) error {
	result := config.DB.Create(run)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
