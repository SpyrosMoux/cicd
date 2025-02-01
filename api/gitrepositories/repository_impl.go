package gitrepositories

import (
	"github.com/spyrosmoux/cicd/api/config"
	"gorm.io/gorm"
)

type gitRepositoryRepository struct {
	db *gorm.DB
}

func NewGitRepositoryRepository(db *gorm.DB) GitRepositoryRepository {
	return &gitRepositoryRepository{db: db}
}

func (r gitRepositoryRepository) SaveGitRepository(gitRepository *GitRepository) error {
	result := config.DB.Create(gitRepository)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r gitRepositoryRepository) FindGitRepositoryById(repoId string) (*GitRepository, error) {
	var repo GitRepository
	queryResult := config.DB.Where("id = ?", repoId).First(&repo)
	if queryResult.Error != nil {
		return nil, queryResult.Error
	}
	return &repo, nil
}

func (r gitRepositoryRepository) UpdateGitRepository(gitRepository *GitRepository) error {
	result := config.DB.Save(gitRepository)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r gitRepositoryRepository) FindAllGitRepository() (*[]GitRepository, error) {
	var gitRepos *[]GitRepository
	result := config.DB.Find(&gitRepos)
	if result.Error != nil {
		return &[]GitRepository{}, result.Error
	}

	return gitRepos, nil
}

func (r gitRepositoryRepository) DeleteGitRepository(repoId string) error {
	result := config.DB.Delete(&GitRepository{}, "id = ?", repoId)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
