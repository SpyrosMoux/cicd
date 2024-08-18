package project

import (
	"log"
)

type ProjectsServiceImpl struct {
	ProjectsRepository ProjectsRepository
}

func NewProjectsServiceImpl(projectsRepository ProjectsRepository) ProjectsService {
	return &ProjectsServiceImpl{ProjectsRepository: projectsRepository}
}

func (p ProjectsServiceImpl) Create(project Project) *Project {
	newProject, err := p.ProjectsRepository.Save(&project)
	if err != nil {
		log.Println(err)
		return nil
	}

	return newProject
}

func (p ProjectsServiceImpl) Update(project Project) {
	//TODO(spyrosmoux) implement me
	panic("implement me")
}

func (p ProjectsServiceImpl) Delete(projectId string) error {
	err := p.ProjectsRepository.Delete(projectId)
	if err != nil {
		return err
	}

	return nil
}

func (p ProjectsServiceImpl) FindById(projectId string) *Project {
	project, err := p.ProjectsRepository.FindById(projectId)
	if err != nil {
		log.Println(err)
		return nil
	}

	return project
}

func (p ProjectsServiceImpl) FindAll() *[]Project {
	project, err := p.ProjectsRepository.FindAll()
	if err != nil {
		log.Println(err)
		return nil
	}

	return project
}
