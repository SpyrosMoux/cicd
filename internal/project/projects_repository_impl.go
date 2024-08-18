package project

import (
	"gorm.io/gorm"
)

type ProjectsRepositoryImpl struct {
	Db *gorm.DB
}

func NewProjectsRepositoryImpl(Db *gorm.DB) ProjectsRepository {
	return &ProjectsRepositoryImpl{Db: Db}
}

func (p ProjectsRepositoryImpl) Save(project *Project) (*Project, error) {
	result := p.Db.Create(&project)
	if result.Error != nil {
		return nil, result.Error
	}

	return project, nil
}

func (p ProjectsRepositoryImpl) Update(project *Project) {
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

func (p ProjectsRepositoryImpl) FindById(projectId string) (*Project, error) {
	project := &Project{}
	result := p.Db.Find(&project, "id = ?", projectId)

	if result.Error != nil {
		return nil, result.Error
	}

	return project, nil
}

func (p ProjectsRepositoryImpl) FindAll() (*[]Project, error) {
	project := &[]Project{}
	result := p.Db.Find(&project)

	if result.Error != nil {
		return nil, result.Error
	}

	return project, nil
}
