package repository

type RepositoriesRepository interface {
	Save(repository *Repository) (*Repository, error)
	Update(repository *Repository) (*Repository, error)
	Delete(repositoryId string) error
	FindById(repositoryId string) (*Repository, error)
	FindAll() ([]*Repository, error)
}
