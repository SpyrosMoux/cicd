package services

import (
	"spyrosmoux/api/internal/models"
	"spyrosmoux/api/internal/repositories"
)

type RepositoriesServiceImpl struct {
	RepositoriesRepository repositories.RepositoriesRepository
}

func NewRepositoriesServiceImpl(repositoriesRepository repositories.RepositoriesRepository) RepositoriesService {
	return &RepositoriesServiceImpl{RepositoriesRepository: repositoriesRepository}
}

func (r RepositoriesServiceImpl) Create(repository models.Repository) *models.Repository {
	newRepository, err := r.RepositoriesRepository.Save(&repository)
	if err != nil {
		return nil
	}

	return newRepository
}

func (r RepositoriesServiceImpl) Update(repository models.Repository) {
	//TODO implement me
	panic("implement me")
}

func (r RepositoriesServiceImpl) Delete(repositoryId string) error {
	err := r.RepositoriesRepository.Delete(repositoryId)
	if err != nil {
		return err
	}

	return nil
}

func (r RepositoriesServiceImpl) FindById(repositoryId string) (*models.Repository, error) {
	repository, err := r.RepositoriesRepository.FindById(repositoryId)
	if err != nil {
		return nil, err
	}

	return repository, nil
}

func (r RepositoriesServiceImpl) FindAll() ([]*models.Repository, error) {
	repositoriesFound, err := r.RepositoriesRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return repositoriesFound, nil
}
