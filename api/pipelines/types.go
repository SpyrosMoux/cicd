package pipelines

import "github.com/google/uuid"

type Pipeline struct {
	Id              string `json:"id" gorm:"primary_key"`
	Name            string `json:"name"`
	GitRepositoryId string `json:"gitrepository_id" gorm:"uniqueIndex:idx_name_repo"`
	Filepath        string `json:"filepath"`
}

type PipelineDto struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	GitRepositoryId string `json:"gitrepository_id"`
	Filepath        string `json:"filepath"`
}

type CreatePipelineDto struct {
	Name            string `json:"name"`
	GitRepositoryId string `json:"gitrepository_id"`
	Filepath        string `json:"filepath"`
}

func NewPipeline(name, gitRepoId, filepath string) *Pipeline {
	return &Pipeline{
		Id:              uuid.New().String(),
		Name:            name,
		GitRepositoryId: gitRepoId,
		Filepath:        filepath,
	}
}

func toPipelineDto(pipeline Pipeline) *PipelineDto {
	return &PipelineDto{
		Id:              pipeline.Id,
		Name:            pipeline.Name,
		GitRepositoryId: pipeline.GitRepositoryId,
		Filepath:        pipeline.Filepath,
	}
}
