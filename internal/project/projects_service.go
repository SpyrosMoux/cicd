package project

type ProjectsService interface {
	Create(project Project) *Project
	Update(project Project)
	Delete(projectId string) error
	FindById(projectId string) *Project
	FindAll() *[]Project
}
