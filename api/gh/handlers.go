package gh

import (
	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/cicd/common/helpers"
)

var WebhookSecret = helpers.LoadEnvVariable("GH_WEBHOOK_SECRET")
var GhToken = helpers.LoadEnvVariable("GH_TOKEN")

type Handler interface {
	HandleWebhook(ctx *gin.Context)
}
