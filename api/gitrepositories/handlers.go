package gitrepositories

import (
	"github.com/gin-gonic/gin"
)

type GitRepositoryHandler interface {
	HandleCreateGitRepository(ctx *gin.Context)
	HandleGetGitRepositoryById(ctx *gin.Context)
	HandleUpdateGitRepository(ctx *gin.Context)
	HandleGetAllGitRepository(ctx *gin.Context)
	HandleDeleteGitRepository(ctx *gin.Context)
}
