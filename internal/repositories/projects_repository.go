package repositories

import (
	"spyrosmoux/api/internal/models"
)

type ProjectsRepository interface {
	Save(project *models.Project) (*models.Project, error)
	Update(project *models.Project)
	Delete(projectId string) error
	FindById(projectId string) (*models.Project, error)
	FindAll() (*[]models.Project, error)
}
