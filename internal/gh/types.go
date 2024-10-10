package gh

import (
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"github.com/spyrosmoux/api/internal/pipelineruns"
	"github.com/spyrosmoux/api/internal/queue"
)

type PushEventAdapter struct {
	Event *github.PushEvent
}

func (eventAdapter *PushEventAdapter) HandleGhEvent() {
	event := eventAdapter.Event

	fmt.Printf("Received a push event for ref %s\n", *event.Ref)

	pipelines, err := FetchPipelineConfig(*event.Repo.Owner.Name, *event.Repo.Name, *event.Ref, *event.Installation.ID)
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

		err := pipelineruns.AddPipelineRun(pipelineRun)
		if err != nil {
			log.Printf("Failed to add pipeline run: %v", err)
			return
		}

		fmt.Println("Publishing pipeline run with id: " + pipelineRun.Id)
		queue.PublishJob(pipelineRun.Id, fmt.Sprintf("%v", pipeline))
	}
}

func matchPushEventWithBranch(event *github.PushEvent, branches []string) bool {
	shouldRun := false
	for _, branch := range branches {
		if event.GetRef() == branch {
			fmt.Printf("Matching push event for branch %s\n", branch)
			shouldRun = true
		}
	}

	return shouldRun
}
