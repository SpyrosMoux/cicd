package gitrepositories

import (
	"github.com/gin-gonic/gin"
)

func Routes(route *gin.Engine, gitRepositoryHandler GitRepositoryHandler) {
	routes := route.Group("/app/cicd/api/repositories")
	{
		routes.POST("", gitRepositoryHandler.HandleCreateGitRepository)
		routes.GET("/:id", gitRepositoryHandler.HandleGetGitRepositoryById)
		routes.PUT("/:id", gitRepositoryHandler.HandleUpdateGitRepository)
		routes.GET("", gitRepositoryHandler.HandleGetAllGitRepository)
		routes.DELETE("/:id", gitRepositoryHandler.HandleDeleteGitRepository)
	}
}
