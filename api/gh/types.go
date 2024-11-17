package gh

import (
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/spyrosmoux/cicd/api/config"
	"github.com/spyrosmoux/cicd/api/entities"
	"github.com/spyrosmoux/cicd/api/repositories"
	"github.com/spyrosmoux/cicd/api/services"
	"github.com/spyrosmoux/cicd/common/queue"
	"gopkg.in/yaml.v3"
	"log"
	"log/slog"
	"strings"
)

type PushEventAdapter struct {
	Event *github.PushEvent
}

func (eventAdapter *PushEventAdapter) HandleGhEvent() {
	// TODO(@spyrosmoux) is there a better way to inject the service?
	repo := repositories.NewPipelineRunsRepository(config.DB)
	svc := services.NewPipelineRunsService(repo)

	event := eventAdapter.Event

	fmt.Printf("Received a push event for ref %s\n", *event.Ref)

	pipelines, err := FetchPipelineConfig(*event.Repo.Owner.Name, *event.Repo.Name, *event.Ref, *event.Installation.ID)
	if err != nil {
		log.Printf("Failed to fetch pipeline config: %v", err)
	}

	// Publish all triggered pipelines
	for _, pipeline := range pipelines {
		pipelineRun := entities.NewPipelineRun(*event.Repo.Name, *event.Ref)

		if !matchPushEventWithBranch(event, pipeline.Triggers.Branch) {
			fmt.Printf("No matching push event for branch %s\n", *event.Ref)
			continue
		}

		err := svc.AddPipelineRun(pipelineRun)
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

// A range of events could trigger a pipeline. i.e push a new branch, make a commit
// For now we only support push events
// These events should be matched as follows
// - If a commit is made to the specified branch -> run
// - If a branch is created, and specified in the triggers -> run
// - If a tag is create, and the tag is specified in the trigges -> run
// Apart from creating stuff the push event also represents deletion events. Such as deleting a tag or branch.
// These events should be ignored
func matchPushEventWithBranch(event *github.PushEvent, branches []string) bool {
	branchName, err := getBranchNameFromRef(event.GetRef())
	if err != nil {
		log.Printf("Failed to get branch name from ref %s: %v", event.GetRef(), err)
		return false
	}

	shouldRun := false
	for _, branch := range branches {
		if branch == "*" {
			shouldRun = true
			break
		}

		if branchName == branch {
			fmt.Printf("Matching push event for branch %s\n", branch)
			shouldRun = true
		}
	}

	return shouldRun
}

func getBranchNameFromRef(ref string) (string, error) {
	if !strings.Contains(ref, "refs/heads/") {
		return "", errors.New("ref does not contain refs/heads/")
	}

	parts := strings.SplitAfter(ref, "refs/heads/")
	if len(parts) != 2 {
		return "", errors.New("error getting branch name from ref " + ref)
	}

	return parts[1], nil
}
