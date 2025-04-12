package logcollector

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/cicd/common/db"
)

func SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:63342"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowCredentials: true,
	}))

	logRepo := NewLogRepository(db.DB)
	logSvc := NewLogService(logRepo)
	logHandler := NewHandler(logSvc)

	apiGroup := router.Group("/app/cicd/logs")
	{
		apiGroup.GET("/health", logHandler.HandleHealth)
		apiGroup.GET("/:runId", logHandler.HandleGetLogsByRunId)
		apiGroup.GET("/ws/:runId", logHandler.HandleStreamLogsByRunId)
	}

	return router
}
