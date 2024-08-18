package repository

type RepositoriesServiceImpl struct {
	RepositoriesRepository RepositoriesRepository
}

func NewRepositoriesServiceImpl(repositoriesRepository RepositoriesRepository) RepositoriesService {
	return &RepositoriesServiceImpl{RepositoriesRepository: repositoriesRepository}
}

func (r RepositoriesServiceImpl) Create(repository Repository) *Repository {
	newRepository, err := r.RepositoriesRepository.Save(&repository)
	if err != nil {
		return nil
	}

	return newRepository
}

func (r RepositoriesServiceImpl) Update(repository Repository) {
	//TODO(spyrosmoux) implement me
	panic("implement me")
}

func (r RepositoriesServiceImpl) Delete(repositoryId string) error {
	err := r.RepositoriesRepository.Delete(repositoryId)
	if err != nil {
		return err
	}

	return nil
}

func (r RepositoriesServiceImpl) FindById(repositoryId string) (*Repository, error) {
	repository, err := r.RepositoriesRepository.FindById(repositoryId)
	if err != nil {
		return nil, err
	}

	return repository, nil
}

func (r RepositoriesServiceImpl) FindAll() ([]*Repository, error) {
	repositoriesFound, err := r.RepositoriesRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return repositoriesFound, nil
}
