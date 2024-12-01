package pipelines

import "github.com/spyrosmoux/cicd/common/dto"

type Service interface {
	PrepareRun(repoMeta dto.Metadata) error
	RunPipeline(pipeline Pipeline, runMetadata dto.Metadata) error
	CleanupRun() error
	ExecuteStep(step Step, variables map[string]string) error
	ExecuteJob(job Job, variables map[string]string) error
	SubstituteVariables(command string, variables map[string]string) string
}
