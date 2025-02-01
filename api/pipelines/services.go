package pipelines

type PipelineService interface {
	CreatePipeline(pipelineDto CreatePipelineDto) (*PipelineDto, error)
	GetAllPipelines() ([]*PipelineDto, error)
}
