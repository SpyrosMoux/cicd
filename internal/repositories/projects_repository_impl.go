package repositories

import (
	"spyrosmoux/api/internal/models"

	"gorm.io/gorm"
)

type ProjectsRepositoryImpl struct {
	Db *gorm.DB
}

func NewProjectsRepositoryImpl(Db *gorm.DB) ProjectsRepository {
	return &ProjectsRepositoryImpl{Db: Db}
}

func (p ProjectsRepositoryImpl) Save(project *models.Project) (*models.Project, error) {
	result := p.Db.Create(&project)
	if result.Error != nil {
		return nil, result.Error
	}

	return project, nil
}

func (p ProjectsRepositoryImpl) Update(project *models.Project) {
	//TODO(spyrosmoux) implement me
	panic("implement me")
}

func (p ProjectsRepositoryImpl) Delete(projectId string) error {
	project, err := p.FindById(projectId)
	if err != nil {
		return err
	}

	result := p.Db.Delete(&project)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (p ProjectsRepositoryImpl) FindById(projectId string) (*models.Project, error) {
	project := &models.Project{}
	result := p.Db.Find(&project, "id = ?", projectId)

	if result.Error != nil {
		return nil, result.Error
	}

	return project, nil
}

func (p ProjectsRepositoryImpl) FindAll() (*[]models.Project, error) {
	project := &[]models.Project{}
	result := p.Db.Find(&project)

	if result.Error != nil {
		return nil, result.Error
	}

	return project, nil
}
