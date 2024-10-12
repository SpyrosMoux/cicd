package handlers

import (
	"log"
	"net/http"

	"github.com/spyrosmoux/api/internal/gh"
	"github.com/spyrosmoux/api/internal/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
)

var GhWebhookSecret = helpers.LoadEnvVariable("GH_WEBHOOK_SECRET")

type GhEventHandler interface {
	HandleGhEvent()
}

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
		handler := gh.PushEventAdapter{
			Event: event,
		}
		handler.HandleGhEvent()
	default:
		log.Printf("Unhandled event type: %s", github.WebHookType(c.Request))
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "bad request",
			"error":  "Unsupported event type",
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
