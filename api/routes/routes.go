package routes

import (
	"github.com/spyrosmoux/cicd/api/config"
	"github.com/spyrosmoux/cicd/api/gh"
	"github.com/spyrosmoux/cicd/api/gitrepositories"
	"github.com/spyrosmoux/cicd/api/pipelineruns"
	"github.com/spyrosmoux/cicd/api/pipelines"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

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

	ghService := gh.NewService(pipelineRunsSvc)
	ghHandler := gh.NewHandler(ghService)
	gh.Routes(router, ghHandler)

	gitRepoRepository := gitrepositories.NewGitRepositoryRepository(config.DB)
	gitRepoSvc := gitrepositories.NewGitRepositoryService(gitRepoRepository)
	gitRepoHandler := gitrepositories.NewGitRepositoryHandler(gitRepoSvc)
	gitrepositories.Routes(router, gitRepoHandler)

	pipelineRepository := pipelines.NewPipelineRepository(config.DB)
	pipelineService := pipelines.NewPipelineService(pipelineRepository, gitRepoSvc)
	pipelineHandler := pipelines.NewPipelineHandler(pipelineService)
	pipelines.Routes(router, pipelineHandler)

	return router
}
