package gitrepositories

type GitRepositoryRepository interface {
	SaveGitRepository(gitRepository *GitRepository) error
	FindGitRepositoryById(repoId string) (*GitRepository, error)
	UpdateGitRepository(gitRepository *GitRepository) error
	FindAllGitRepository() (*[]GitRepository, error)
	DeleteGitRepository(repoId string) error
}
