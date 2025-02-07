package gh

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v68/github"
	"github.com/sirupsen/logrus"
)

type handler struct {
	svc    Service
	logger *logrus.Logger
}

func NewHandler(svc Service, logger *logrus.Logger) Handler {
	return &handler{
		svc:    svc,
		logger: logger,
	}
}

func (h *handler) HandleWebhook(ctx *gin.Context) {
	payload, err := github.ValidatePayload(ctx.Request, []byte(WebhookSecret))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		h.logger.WithError(err).Error("invalid payload")
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(ctx.Request), payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse webhook"})
		h.logger.WithError(err).Error("failed to parse webhook")
		return
	}

	err = h.svc.ProcessEvent(event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.logger.WithError(err).Error("failed to process event")
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
