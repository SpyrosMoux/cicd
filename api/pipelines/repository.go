package pipelines

type PipelineRepository interface {
	SavePipeline(pipeline *Pipeline) error
	FindAllPipelines() ([]Pipeline, error)
}
