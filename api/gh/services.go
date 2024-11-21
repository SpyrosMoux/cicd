package gh

import (
	"github.com/google/go-github/github"
	"github.com/spyrosmoux/cicd/runner/pipelines"
)

type Service interface {
	FetchValidPipelines(repoOwner string, repoName string, branchName string, installationId int64) ([]pipelines.Pipeline, error)
	ProcessEvent(event interface{}) error
	ProcessPushEvent(event *github.PushEvent) error
}
