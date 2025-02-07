package routes

import (
	"github.com/spyrosmoux/cicd/api/config"
	"github.com/spyrosmoux/cicd/api/gh"
	"github.com/spyrosmoux/cicd/api/pipelineruns"
	"github.com/spyrosmoux/cicd/common/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())

	/* SuperTokens Routers */

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:63342"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowCredentials: true,
	}))

	pipelineRunsRepo := pipelineruns.NewRepository(config.DB)
	pipelineRunsSvc := pipelineruns.NewService(pipelineRunsRepo)
	pipelineRunsHandler := pipelineruns.NewHandler(pipelineRunsSvc)
	pipelineruns.Routes(router, pipelineRunsHandler)

	logger := logger.NewLogger()

	ghService := gh.NewService(pipelineRunsSvc, logger)
	ghHandler := gh.NewHandler(ghService, logger)
	gh.Routes(router, ghHandler)

	return router
}
