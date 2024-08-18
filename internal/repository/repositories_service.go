package repository

type RepositoriesService interface {
	Create(repository Repository) *Repository
	Update(repository Repository)
	Delete(repositoryId string) error
	FindById(repositoryId string) (*Repository, error)
	FindAll() ([]*Repository, error)
}
