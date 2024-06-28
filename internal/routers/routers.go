package routers

import (
	"github.com/gin-gonic/gin"
	"spyrosmoux/api/internal/handlers"
)

func SetupRouter() *gin.Engine {
  router := gin.Default()

	router.POST("/webhook", handlers.HandleWebhook)

	return router
}
