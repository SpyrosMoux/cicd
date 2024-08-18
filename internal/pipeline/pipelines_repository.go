package pipeline

type PipelinesRepository interface {
	Save(pipeline *Pipeline) (*Pipeline, error)
	Update(pipeline *Pipeline)
	Delete(pipelineId string) error
	FindById(pipelineId string) (*Pipeline, error)
	FindAll() (*[]Pipeline, error)
}
