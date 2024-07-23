package services

import "spyrosmoux/api/internal/models"

type RepositoriesService interface {
	Create(repository models.Repository) *models.Repository
	Update(repository models.Repository)
	Delete(repositoryId string) error
	FindById(repositoryId string) (*models.Repository, error)
	FindAll() ([]*models.Repository, error)
}
