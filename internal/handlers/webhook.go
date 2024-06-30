package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/webhooks/v6/github"
	"log"
)

func HandleWebhook(c *gin.Context) {
	hook, _ := github.New(github.Options.Secret("MyGitHubSuperSecretSecret...?"))

	payload, err := hook.Parse(c.Request, github.ReleaseEvent, github.PullRequestEvent)
	if err != nil {
		log.Panicf("Error parsing webhook payload: %v", err)
	}

	log.Printf("Received webhook payload: %+v", payload)
}
