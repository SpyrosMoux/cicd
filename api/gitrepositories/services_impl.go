package gitrepositories

import (
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type gitRepositoryService struct {
	gitRepositoryRepository GitRepositoryRepository
}

func NewGitRepositoryService(gitRepo GitRepositoryRepository) GitRepositoryService {
	return &gitRepositoryService{gitRepositoryRepository: gitRepo}
}

func (gitRepoSvc gitRepositoryService) CreateGitRepository(createGitRepositoryDto CreateGitRepositoryDto) (*GitRepositoryDto, error) {
	gitRepo := NewGitRepository(
		createGitRepositoryDto.Name,
		createGitRepositoryDto.Owner,
		createGitRepositoryDto.OwnerEmail,
		createGitRepositoryDto.Url,
	)

	err := gitRepoSvc.gitRepositoryRepository.SaveGitRepository(gitRepo)
	if err != nil {
		slog.Error("unable to add gitrepository", "gitrepository", createGitRepositoryDto.Name, "err", err)
		return &GitRepositoryDto{}, err
	}

	slog.Info("created gitrepository", "gitrepository", gitRepo.Name)
	return toGitRepositoryDto(*gitRepo), nil
}

func (gitRepoSvc gitRepositoryService) findGitRepositoryById(gitRepoId string) (*GitRepository, error) {
	gitRepo, err := gitRepoSvc.gitRepositoryRepository.FindGitRepositoryById(gitRepoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("gitrepository not found", "id", gitRepoId)
			return nil, err
		}
		slog.Error("failure while searching for gitrepository", "id", gitRepoId, "err", err)
		return nil, err
	}
	slog.Info("gitrepository found", "id", gitRepoId)
	return gitRepo, nil
}

func (gitRepoSvc gitRepositoryService) GetGitRepositoryById(gitRepoId string) (*GitRepositoryDto, error) {
	gitRepo, err := gitRepoSvc.findGitRepositoryById(gitRepoId)
	if err != nil {
		return &GitRepositoryDto{}, err
	}
	return toGitRepositoryDto(*gitRepo), nil
}

func (gitRepoSvc gitRepositoryService) UpdateGitRepository(gitRepoId string, gitRepoDto UpdateGitRepositoryDto) (*GitRepositoryDto, error) {
	gitRepo, err := gitRepoSvc.findGitRepositoryById(gitRepoId)
	if err != nil {
		return &GitRepositoryDto{}, err
	}

	gitRepo.Name = gitRepoDto.Name
	gitRepo.FullName = gitRepoDto.Owner + "/" + gitRepoDto.Name
	gitRepo.Owner = gitRepoDto.Owner
	gitRepo.OwnerEmail = gitRepoDto.OwnerEmail
	gitRepo.Url = gitRepoDto.Url

	err = gitRepoSvc.gitRepositoryRepository.UpdateGitRepository(gitRepo)
	if err != nil {
		slog.Error("unable to update gitrepository", "id", gitRepoId, "err", err)
		return &GitRepositoryDto{}, err
	}

	return toGitRepositoryDto(*gitRepo), nil
}

func (gitRepoSvc gitRepositoryService) GetAllGitRepository() (*[]GitRepositoryDto, error) {
	gitRepos, err := gitRepoSvc.gitRepositoryRepository.FindAllGitRepository()
	if err != nil {
		slog.Error("unable to fetch all gitrepositories", "err", err)
		return &[]GitRepositoryDto{}, err
	}

	var gitRepoDtos []GitRepositoryDto
	for _, gitRepo := range *gitRepos {
		gitRepoDtos = append(gitRepoDtos, *toGitRepositoryDto(gitRepo))
	}
	return &gitRepoDtos, nil
}

func (gitRepoSvc gitRepositoryService) DeleteGitRepository(repoId string) error {
	_, err := gitRepoSvc.GetGitRepositoryById(repoId)
	if err != nil {
		return err
	}

	err = gitRepoSvc.gitRepositoryRepository.DeleteGitRepository(repoId)
	if err != nil {
		slog.Error("failed to delete gitrepository", "id", repoId, "err", err)
		return err
	}

	slog.Info("gitrepository deleted", "id", repoId)
	return nil
}
