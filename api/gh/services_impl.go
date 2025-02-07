package gh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/go-github/v68/github"
	"github.com/sirupsen/logrus"
	"github.com/spyrosmoux/cicd/api/pipelineruns"
	"github.com/spyrosmoux/cicd/common/dto"
	"github.com/spyrosmoux/cicd/common/queue"
	"github.com/spyrosmoux/cicd/runner/pipelines"
	"gopkg.in/yaml.v3"
)

type service struct {
	pipelineRunsService pipelineruns.Service
	logger              *logrus.Logger
}

func NewService(pipelineRunsService pipelineruns.Service, logger *logrus.Logger) Service {
	return &service{
		pipelineRunsService: pipelineRunsService,
		logger:              logger,
	}
}

func (svc *service) FetchValidPipelines(repoOwner string, repoName string, branchName string) ([]pipelines.Pipeline, error) {
	ctx := context.Background()
	client := github.NewClient(nil).WithAuthToken(GhToken)
	options := &github.RepositoryContentGetOptions{
		Ref: branchName,
	}

	_, contents, _, err := client.Repositories.GetContents(ctx, repoOwner, repoName, ".flowforge", options) // Only look in .flowforge dir for pipelines
	if err != nil {
		return nil, err
	}

	var validYAMLs []pipelines.Pipeline

	for _, file := range contents {
		svc.logger.WithFields(logrus.Fields{
			"file": file.GetName(),
		}).Info("found file")

		downloadURL := file.GetDownloadURL()
		if downloadURL == "" {
			return nil, errors.New("Could not get download URL for pipeline " + file.GetName())
		}
		fileContent, err := downloadYAMLContent(downloadURL)
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

func (svc *service) ProcessEvent(event interface{}) error {
	switch ghEvent := event.(type) {
	case *github.PushEvent:
		return svc.ProcessPushEvent(ghEvent)
	case *github.PullRequestEvent:
		return svc.ProcessPullRequestEvent(ghEvent)
	default:
		return fmt.Errorf("unsupported event type %T", event)
	}
}

func (svc *service) ProcessPushEvent(event *github.PushEvent) error {
	svc.logger.WithFields(logrus.Fields{
		"ref": *event.Ref,
	}).Info("received a push event for")

	validPipelines, err := svc.FetchValidPipelines(*event.Repo.Owner.Name, *event.Repo.Name, *event.Ref)
	if err != nil {
		svc.logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to fetch pipeline config")
		return err
	}

	// Publish all triggered validPipelines
	for _, pipeline := range validPipelines {
		pipelineRun := pipelineruns.NewPipelineRun(*event.Repo.Name, *event.Ref, "", *event.Sender.Login, pipelineruns.PUSH)

		if !matchPushEventWithBranch(event, pipeline.Triggers.Branch) {
			fmt.Printf("No matching push event for branch %s\n", *event.Ref)
			svc.logger.WithFields(logrus.Fields{
				"ref": *event.Ref,
			}).Info("no matching push event for branch")
			continue
		}

		response := svc.pipelineRunsService.AddPipelineRun(pipelineRun)
		if response.Error != "" {
			err := errors.New(response.Error)
			svc.logger.WithError(err).Error("failed to add pipeline")
			return err
		}

		pipelineAsBytes, err := yaml.Marshal(pipeline)
		if err != nil {
			svc.logger.WithError(err).Error("unable to convert pipeline yaml to string")
			return err
		}

		repoVisibility := dto.PUBLIC
		if *event.Repo.Private {
			repoVisibility = dto.PRIVATE
		}

		branch, err := getBranchNameFromRef(*event.Ref)
		if err != nil {
			return err
		}
		publishRunDto := dto.PublishRunDto{
			PipelineAsBytes: pipelineAsBytes,
			Metadata: dto.Metadata{
				Repository:     *event.Repo.Name,
				Branch:         branch,
				RepoOwner:      *event.Repo.Owner.Name,
				RepoVisibility: repoVisibility,
				VcsSource:      dto.GITHUB,
				VcsToken:       GhToken,
			},
		}

		publishBody, err := json.Marshal(publishRunDto)
		if err != nil {
			svc.logger.WithError(err).Error("unable to marshal publishRunDto into Json")
			return err
		}

		svc.logger.WithFields(logrus.Fields{
			"pipelineRunId": pipelineRun.Id,
		}).Info("publishing pipeline run")
		queue.PublishJob(pipelineRun.Id, publishBody)
	}

	return nil
}

// ProcessPullRequestEvent will run if the destination branch of the pull request
// is in the list of pr: in the yaml
func (svc *service) ProcessPullRequestEvent(event *github.PullRequestEvent) error {
	headBranch := event.GetPullRequest().GetHead()
	baseBranch := event.GetPullRequest().GetBase()
	svc.logger.WithFields(logrus.Fields{
		"repo":    event.GetRepo().GetName(),
		"headRef": headBranch.GetRef(),
		"baseRef": baseBranch.GetRef(),
	}).Info("received a PullRequestEvent")

	switch event.GetAction() {
	case "opened", "reopened", "synchronize":
		svc.logger.WithFields(logrus.Fields{
			"action": event.GetAction(),
		}).Debug("valid")
	default:
		svc.logger.WithFields(logrus.Fields{
			"action": event.GetAction(),
		}).Warn("skipping pull request event")
		return nil
	}

	validPipelines, err := svc.FetchValidPipelines(headBranch.GetRepo().GetOwner().GetLogin(), headBranch.GetRepo().GetName(), headBranch.GetRef())
	if err != nil {
		svc.logger.WithFields(logrus.Fields{
			"repo": headBranch.GetRepo().GetName(),
			"err":  err,
		}).Error("unable to fetch valid pipelines")
		return err
	}

	for _, pipeline := range validPipelines {
		pipelineRun := pipelineruns.NewPipelineRun(headBranch.GetRepo().GetName(), headBranch.GetRef(), "", event.GetSender().GetLogin(), pipelineruns.PR)

		if !matchPullRequestEventWithBranch(event, pipeline.Triggers.PR) {
			svc.logger.WithFields(logrus.Fields{
				"branch": baseBranch.GetRef(),
			}).Info("no matching base")
			continue
		}

		response := svc.pipelineRunsService.AddPipelineRun(pipelineRun)
		if response.Error != "" {
			svc.logger.WithError(err).Error("failed to add pipelineRun")
			return err
		}

		pipelineAsBytes, err := yaml.Marshal(pipeline)
		if err != nil {
			svc.logger.WithError(err).Error("unable to convert pipeline yaml to string")
			return err
		}

		repoVisibility := dto.PUBLIC
		if event.GetRepo().GetPrivate() {
			repoVisibility = dto.PRIVATE
		}

		branch, err := getBranchNameFromRef(headBranch.GetRef())
		if err != nil {
			return err
		}
		publishRunDto := dto.PublishRunDto{
			PipelineAsBytes: pipelineAsBytes,
			Metadata: dto.Metadata{
				Repository:     event.GetRepo().GetName(),
				Branch:         branch,
				RepoOwner:      event.GetRepo().GetOwner().GetLogin(),
				RepoVisibility: repoVisibility,
				VcsSource:      dto.GITHUB,
				VcsToken:       GhToken,
			},
		}

		publishBody, err := json.Marshal(publishRunDto)
		if err != nil {
			svc.logger.WithError(err).Error("unable to marshal pipelineRunDto into Json")
			return err
		}

		svc.logger.WithFields(logrus.Fields{
			"id": pipelineRun.Id,
		}).Info("publishing pipeline run")
		queue.PublishJob(pipelineRun.Id, publishBody)

	}
	return nil
}
