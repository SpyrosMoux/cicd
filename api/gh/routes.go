package gh

import (
	"github.com/gin-gonic/gin"
)

func Routes(rg *gin.RouterGroup, ghHandler Handler) {
	ghGroup := rg.Group("/gh")
	{
		ghGroup.POST("/webhook", ghHandler.HandleWebhook)
	}
}
