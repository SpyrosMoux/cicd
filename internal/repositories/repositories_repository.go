package repositories

import "spyrosmoux/api/internal/models"

type RepositoriesRepository interface {
	Save(repository *models.Repository) (*models.Repository, error)
	Update(repository *models.Repository) (*models.Repository, error)
	Delete(repositoryId string) error
	FindById(repositoryId string) (*models.Repository, error)
	FindAll() ([]*models.Repository, error)
}
