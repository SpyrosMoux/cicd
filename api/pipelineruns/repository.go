package pipelineruns

type Repository interface {
	FindById(pipelineRunId string) (*PipelineRun, error)
	FindAll() (*[]PipelineRun, error)
	Update(run *PipelineRun) (*PipelineRun, error)
	Create(run *PipelineRun) error
}
