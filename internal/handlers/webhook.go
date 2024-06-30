package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/webhooks/v6/github"
	"io/ioutil"
	"log"
	"net/http"
	"spyrosmoux/api/internal/helpers"
)

var GhSecret = helpers.LoadEnvVariable("GH_SECRET")
var GhToken = helpers.LoadEnvVariable("GH_TOKEN")

func HandleWebhook(c *gin.Context) {
	hook, _ := github.New(github.Options.Secret(GhSecret))

	payload, err := hook.Parse(c.Request, github.PushEvent)
	if err != nil {
		log.Panicf("Error parsing webhook payload: %v", err)
	}

	push := payload.(github.PushPayload)

	pipeline, err := fetchPipelineConfig(push.Repository.FullName, push.Ref)
	if err != nil {
		log.Panicf("Error fetching pipeline config: %v", err)
	}

	log.Printf("Received pipeline config: %+v", pipeline)
}

func fetchPipelineConfig(repoFullName, branchName string) ([]byte, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/sample-pipeline.yaml?ref=%s", repoFullName, branchName)
	log.Printf("Fetching pipeline config from %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3.raw")
	req.Header.Set("Authorization", "Bearer "+GhToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pipeline config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch pipeline config: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
