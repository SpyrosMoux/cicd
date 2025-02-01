package pipelines

import (
	"log/slog"

	"github.com/spyrosmoux/cicd/api/gitrepositories"
)

type pipelineService struct {
	pipelineRepository   PipelineRepository
	gitRepositoryService gitrepositories.GitRepositoryService
}

func NewPipelineService(pipelineRepository PipelineRepository, gitRepositorySvc gitrepositories.GitRepositoryService) PipelineService {
	return &pipelineService{
		pipelineRepository:   pipelineRepository,
		gitRepositoryService: gitRepositorySvc,
	}
}

func (pipelineService pipelineService) CreatePipeline(pipelineDto CreatePipelineDto) (*PipelineDto, error) {
	_, err := pipelineService.gitRepositoryService.GetGitRepositoryById(pipelineDto.GitRepositoryId)
	if err != nil {
		return &PipelineDto{}, err
	}

	pipeline := NewPipeline(
		pipelineDto.Name,
		pipelineDto.GitRepositoryId,
		pipelineDto.Filepath,
	)

	err = pipelineService.pipelineRepository.SavePipeline(pipeline)
	if err != nil {
		slog.Error("failed to save pipeline", "name", pipelineDto.Name, "gitRepository", pipelineDto.GitRepositoryId, "err", err)
		return &PipelineDto{}, err
	}

	return toPipelineDto(*pipeline), nil
}

func (pipelineService pipelineService) GetAllPipelines() ([]*PipelineDto, error) {
	pipelines, err := pipelineService.pipelineRepository.FindAllPipelines()
	if err != nil {
		slog.Error("failed to get all pipelines", "err", err)
		return nil, err
	}

	var pipelineDtos []*PipelineDto
	for _, pipeline := range pipelines {
		pipelineDtos = append(pipelineDtos, toPipelineDto(pipeline))
	}
	return pipelineDtos, nil
}
