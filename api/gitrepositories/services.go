package gitrepositories

type GitRepositoryService interface {
	CreateGitRepository(createGitRepositoryDto CreateGitRepositoryDto) (*GitRepositoryDto, error)
	findGitRepositoryById(repoId string) (*GitRepository, error)
	GetGitRepositoryById(repoId string) (*GitRepositoryDto, error)
	UpdateGitRepository(repoId string, gitRepoDto UpdateGitRepositoryDto) (*GitRepositoryDto, error)
	GetAllGitRepository() (*[]GitRepositoryDto, error)
	DeleteGitRepository(repoId string) error
}
