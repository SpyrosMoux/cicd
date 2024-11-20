package gh

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/spyrosmoux/cicd/api/config"
	"github.com/spyrosmoux/cicd/api/pipelineruns"
	"github.com/spyrosmoux/cicd/common/queue"
	"github.com/spyrosmoux/cicd/runner/pipelines"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
	"log"
	"log/slog"
)

type Service interface {
	FetchValidPipelines(repoOwner string, repoName string, branchName string, installationId int64) ([]pipelines.Pipeline, error)
	HandlePushEvent(event *github.PushEvent)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (svc *service) FetchValidPipelines(repoOwner string, repoName string, branchName string, installationId int64) ([]pipelines.Pipeline, error) {
	token, err := getInstallationToken(installationId)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{
		TokenType:   "Bearer",
		AccessToken: token,
	})
	tokenClient := oauth2.NewClient(ctx, tokenSource)
	client := github.NewClient(tokenClient)
	options := &github.RepositoryContentGetOptions{
		Ref: branchName,
	}

	_, contents, _, err := client.Repositories.GetContents(ctx, repoOwner, repoName, ".flowforge", options) // Only look in .flowforge dir for pipelines
	if err != nil {
		return nil, err
	}

	var validYAMLs []pipelines.Pipeline

	for _, file := range contents {
		fmt.Printf("Found file: %s\n", file.GetName())

		downloadURL := file.GetDownloadURL()
		if downloadURL == "" {
			return nil, errors.New("Could not get download URL for pipeline " + file.GetName())
		}
		fileContent, err := downloadYAMLContent(downloadURL, installationId)
		if err != nil {
			return nil, err
		}

		pipeline, err := pipelines.ValidateYAMLStructure([]byte(fileContent))
		if err != nil {
			// TODO(@SpyrosMoux) what happens if a yaml is invalid, and there are multiple yamls in the repo?
			// skipping invalid yaml for now
			fmt.Println(err.Error())
			continue // skip to next pipeline
		}

		validYAMLs = append(validYAMLs, pipeline)
	}

	return validYAMLs, nil
}

func (svc *service) HandlePushEvent(event *github.PushEvent) {
	// TODO(@SpyrosMoux) is there a better way to inject the service?
	pipelineRunsRepo := pipelineruns.NewRepository(config.DB)
	pipelineRunsService := pipelineruns.NewService(pipelineRunsRepo)

	fmt.Printf("Received a push event for ref %s\n", *event.Ref)

	pipelines, err := svc.FetchValidPipelines(*event.Repo.Owner.Name, *event.Repo.Name, *event.Ref, *event.Installation.ID)
	if err != nil {
		log.Printf("Failed to fetch pipeline config: %v", err)
	}

	// Publish all triggered pipelines
	for _, pipeline := range pipelines {
		pipelineRun := pipelineruns.NewPipelineRun(*event.Repo.Name, *event.Ref)

		if !matchPushEventWithBranch(event, pipeline.Triggers.Branch) {
			fmt.Printf("No matching push event for branch %s\n", *event.Ref)
			continue
		}

		err := pipelineRunsService.AddPipelineRun(pipelineRun)
		if err != nil {
			log.Printf("Failed to add pipeline run: %v", err)
			return
		}

		pipelineAsString, err := yaml.Marshal(pipeline)
		if err != nil {
			slog.Error("Unable to convert pipeline yaml to string")
			return
		}

		fmt.Println("Publishing pipeline run with id: " + pipelineRun.Id)
		queue.PublishJob(pipelineRun.Id, pipelineAsString)
	}
}
