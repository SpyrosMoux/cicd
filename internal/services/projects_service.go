package services

import "spyrosmoux/api/internal/models"

type ProjectsService interface {
	Create(project models.Project) *models.Project
	Update(project models.Project)
	Delete(projectId string) error
	FindById(projectId string) *models.Project
	FindAll() *[]models.Project
}
