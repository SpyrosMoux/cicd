package gh

import (
	"github.com/gin-gonic/gin"
)

func Routes(route *gin.Engine, ghHandler Handler) {
	routes := route.Group("/app/cicd/api/gh")
	{
		routes.POST("/webhook", ghHandler.HandleWebhook)
	}
}
