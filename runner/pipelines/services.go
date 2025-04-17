package pipelines

import (
	"context"

	"github.com/spyrosmoux/cicd/common/dto"
)

type Service interface {
	PrepareRun(repoMeta dto.Metadata) error
	RunPipeline(ctx context.Context, pipline Pipeline, runMetadata dto.Metadata) error
	CleanupRun() error
	ExecuteStep(ctx context.Context, step Step, variables map[string]string) error
	ExecuteJob(ctx context.Context, job Job, variables map[string]string) error
	SubstituteUserVariables(command string, variables map[string]string) string
	SubstitutePredefinedVariables(pipeline Pipeline, variables map[string]string) Pipeline
}
