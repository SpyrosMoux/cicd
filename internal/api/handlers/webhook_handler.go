package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spyrosmoux/api/internal/pipelineruns"

	"github.com/spyrosmoux/api/internal/gh"
	"github.com/spyrosmoux/api/internal/helpers"
	"github.com/spyrosmoux/api/internal/queue"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
)

var GhWebhookSecret = helpers.LoadEnvVariable("GH_WEBHOOK_SECRET")

func HandleWebhook(c *gin.Context) {
	payload, err := github.ValidatePayload(c.Request, []byte(GhWebhookSecret))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(c.Request), payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse webhook"})
		return
	}

	switch event := event.(type) {
	case *github.PushEvent:
		handlePushEvent(event)
	default:
		log.Printf("Unhandled event type: %s", github.WebHookType(c.Request))
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func handlePushEvent(event *github.PushEvent) {
	fmt.Printf("Received a push event for ref %s\n", *event.Ref)

	pipelines, err := gh.FetchPipelineConfig(*event.Repo.Owner.Name, *event.Repo.Name, *event.Ref, *event.Installation.ID)
	if err != nil {
		log.Printf("Failed to fetch pipeline config: %v", err)
	}

	// Publish all triggered pipelines
	for _, pipeline := range pipelines {
		pipelineRun := pipelineruns.NewPipelineRun(*event.Repo.Name, *event.Ref)

		err := pipelineruns.AddPipelineRun(pipelineRun)
		if err != nil {
			log.Printf("Failed to add pipeline run: %v", err)
			return
		}

		fmt.Println("Publishing pipeline run with id: " + pipelineRun.Id)
		queue.PublishJob(pipelineRun.Id, pipeline)
	}
}
