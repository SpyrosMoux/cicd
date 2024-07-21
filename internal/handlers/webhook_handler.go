package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"io"
	"log"
	"net/http"
	"spyrosmoux/api/internal/auth"
	"spyrosmoux/api/internal/helpers"
	"spyrosmoux/api/internal/queue"
)

var GhWebhookSecret = helpers.LoadEnvVariable("GH_WEBHOOK_SECRET")
var GhToken = helpers.LoadEnvVariable("GH_TOKEN")

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

func fetchPipelineConfig(repoFullName string, branchName string) ([]byte, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/sample-pipeline.yaml?ref=%s", repoFullName, branchName)
	log.Printf("Fetching pipeline config from %s", url)

	token := auth.GenerateJWT()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3.raw")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pipeline config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch pipeline config: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func handlePushEvent(event *github.PushEvent) {
	fmt.Printf("Received a push event for ref %s\n", *event.Ref)

	pipeline, err := fetchPipelineConfig(*event.Repo.FullName, *event.Ref)
	if err != nil {
		log.Printf("Failed to fetch pipeline config: %v", err)
	}

	queue.PublishJob(string(pipeline))
}
