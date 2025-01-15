package gh

import (
	"github.com/google/go-github/v68/github"
	"github.com/spyrosmoux/cicd/runner/pipelines"
)

type Service interface {
	FetchValidPipelines(repoOwner string, repoName string, branchName string) ([]pipelines.Pipeline, error)
	ProcessEvent(event interface{}) error
	ProcessPushEvent(event *github.PushEvent) error
	ProcessPullRequestEvent(event *github.PullRequestEvent) error
}
