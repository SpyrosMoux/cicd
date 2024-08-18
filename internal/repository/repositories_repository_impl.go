package repository

import (
	"gorm.io/gorm"
)

type RepositoriesRepositoryImpl struct {
	Db *gorm.DB
}

func NewRepositoriesRepositoryImpl(Db *gorm.DB) RepositoriesRepository {
	return &RepositoriesRepositoryImpl{Db: Db}
}

func (r RepositoriesRepositoryImpl) Save(repository *Repository) (*Repository, error) {
	result := r.Db.Create(&repository)
	if result.Error != nil {
		return nil, result.Error
	}

	return repository, nil
}

func (r RepositoriesRepositoryImpl) Update(repository *Repository) (*Repository, error) {
	//TODO(spyrosmoux) implement me
	panic("implement me")
}

func (r RepositoriesRepositoryImpl) Delete(repositoryId string) error {
	repository, err := r.FindById(repositoryId)
	if err != nil {
		return err
	}

	result := r.Db.Delete(&repository)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r RepositoriesRepositoryImpl) FindById(repositoryId string) (*Repository, error) {
	repository := &Repository{}
	result := r.Db.Find(&repository, "id = ?", repositoryId)
	if result.Error != nil {
		return nil, result.Error
	}

	return repository, nil
}

func (r RepositoriesRepositoryImpl) FindAll() ([]*Repository, error) {
	repositories := []*Repository{}
	result := r.Db.Find(&repositories)
	if result.Error != nil {
		return nil, result.Error
	}

	return repositories, nil
}
