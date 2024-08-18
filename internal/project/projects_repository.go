package project

type ProjectsRepository interface {
	Save(project *Project) (*Project, error)
	Update(project *Project)
	Delete(projectId string) error
	FindById(projectId string) (*Project, error)
	FindAll() (*[]Project, error)
}
