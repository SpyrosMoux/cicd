package gh

import (
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"log"
	"net/http"
)

type handler struct {
	svc Service
}

func NewHandler(svc Service) Handler {
	return &handler{svc: svc}
}

func (h *handler) HandleWebhook(ctx *gin.Context) {
	payload, err := github.ValidatePayload(ctx.Request, []byte(WebhookSecret))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(ctx.Request), payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse webhook"})
		return
	}

	switch ghEvent := event.(type) {
	case *github.PushEvent:
		h.svc.HandlePushEvent(ghEvent)
	default:
		log.Printf("Unhandled event type: %s", github.WebHookType(ctx.Request))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "bad request",
			"error":  "Unsupported event type",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
