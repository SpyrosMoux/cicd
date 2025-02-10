package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/cicd/api/config"
	"github.com/spyrosmoux/cicd/api/gh"
	"github.com/spyrosmoux/cicd/api/pipelineruns"
	"github.com/spyrosmoux/cicd/common/dto"
	"github.com/spyrosmoux/cicd/common/logger"
	"net/http"
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

	logs := logger.NewLogger()

	ghService := gh.NewService(pipelineRunsSvc, logs)
	ghHandler := gh.NewHandler(ghService, logs)

	apiGroup := router.Group("/app/cicd/api")
	{
		gh.Routes(apiGroup, ghHandler)
		pipelineruns.Routes(apiGroup, pipelineRunsHandler)
		apiGroup.GET("/health", handleHealth)
	}

	return router
}

func handleHealth(ctx *gin.Context) {
	response := dto.NewResponseDto(http.StatusOK, "I'm Alive!", "", "")
	ctx.JSON(response.Status, response)
}
