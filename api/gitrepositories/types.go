package gitrepositories

import (
	"github.com/google/uuid"
)

type GitRepository struct {
	Id         string `json:"id" gorm:"primary_key"`
	Name       string `json:"name"`
	FullName   string `json:"full_name" gorm:"unique"`
	Owner      string `json:"owner"`
	OwnerEmail string `json:"owner_email"`
	Url        string `json:"url" gorm:"unique"`
}

type CreateGitRepositoryDto struct {
	Name       string `json:"name"`
	Owner      string `json:"owner"`
	OwnerEmail string `json:"owner_email"`
	Url        string `json:"url"`
}

type UpdateGitRepositoryDto struct {
	Name       string `json:"name"`
	Owner      string `json:"owner"`
	OwnerEmail string `json:"owner_email"`
	Url        string `json:"url"`
}

type GitRepositoryDto struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	FullName   string `json:"full_name"`
	Owner      string `json:"owner"`
	OwnerEmail string `json:"owner_email"`
	Url        string `json:"url"`
}

func NewGitRepository(name, owner, ownerEmail, url string) *GitRepository {
	return &GitRepository{
		Id:         uuid.New().String(),
		Name:       name,
		FullName:   owner + "/" + name,
		Owner:      owner,
		OwnerEmail: ownerEmail,
		Url:        url,
	}
}

func toGitRepositoryDto(gitRepo GitRepository) *GitRepositoryDto {
	return &GitRepositoryDto{
		Id:         gitRepo.Id,
		Name:       gitRepo.Name,
		FullName:   gitRepo.FullName,
		Owner:      gitRepo.Owner,
		OwnerEmail: gitRepo.OwnerEmail,
		Url:        gitRepo.Url,
	}
}
