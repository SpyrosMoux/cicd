package services

import (
	"log"
	"spyrosmoux/api/internal/models"
	"spyrosmoux/api/internal/repositories"
)

type ProjectsServiceImpl struct {
	ProjectsRepository repositories.ProjectsRepository
}

func NewProjectsServiceImpl(projectsRepository repositories.ProjectsRepository) ProjectsService {
	return &ProjectsServiceImpl{ProjectsRepository: projectsRepository}
}

func (p ProjectsServiceImpl) Create(project models.Project) *models.Project {
	newProject, err := p.ProjectsRepository.Save(&project)
	if err != nil {
		log.Println(err)
		return nil
	}

	return newProject
}

func (p ProjectsServiceImpl) Update(project models.Project) {
	//TODO implement me
	panic("implement me")
}

func (p ProjectsServiceImpl) Delete(projectId string) error {
	err := p.ProjectsRepository.Delete(projectId)
	if err != nil {
		return err
	}

	return nil
}

func (p ProjectsServiceImpl) FindById(projectId string) *models.Project {
	project, err := p.ProjectsRepository.FindById(projectId)
	if err != nil {
		log.Println(err)
		return nil
	}

	return project
}

func (p ProjectsServiceImpl) FindAll() *[]models.Project {
	project, err := p.ProjectsRepository.FindAll()
	if err != nil {
		log.Println(err)
		return nil
	}

	return project
}
