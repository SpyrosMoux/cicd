package gh

import (
	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/cicd/common/helpers"
)

var WebhookSecret = helpers.LoadEnvVariable("GH_WEBHOOK_SECRET")

type Handler interface {
	HandleWebhook(ctx *gin.Context)
}
