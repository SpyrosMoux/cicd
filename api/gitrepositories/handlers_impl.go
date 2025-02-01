package gitrepositories

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/cicd/common/dto"
)

type gitRepositoryHandler struct {
	gitRepoSvc GitRepositoryService
}

func NewGitRepositoryHandler(gitRepoSvc GitRepositoryService) GitRepositoryHandler {
	return &gitRepositoryHandler{gitRepoSvc: gitRepoSvc}
}

func (h gitRepositoryHandler) HandleCreateGitRepository(ctx *gin.Context) {
	var createGitRepositoryDto CreateGitRepositoryDto
	err := ctx.ShouldBindJSON(&createGitRepositoryDto)
	if err != nil {
		response := dto.NewResponseDto(http.StatusBadRequest, "", err.Error(), "")
		ctx.AbortWithStatusJSON(response.Status, response)
		return
	}

	gitRepository, err := h.gitRepoSvc.CreateGitRepository(createGitRepositoryDto)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gitRepository)
}

func (h gitRepositoryHandler) HandleGetGitRepositoryById(ctx *gin.Context) {
	repoId := ctx.Param("id")
	repo, err := h.gitRepoSvc.GetGitRepositoryById(repoId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, repo)
}

func (h gitRepositoryHandler) HandleUpdateGitRepository(ctx *gin.Context) {
	repoId := ctx.Param("id")
	var updateGitRepoDto UpdateGitRepositoryDto
	err := ctx.ShouldBindJSON(&updateGitRepoDto)
	if err != nil {
		response := dto.NewResponseDto(http.StatusBadRequest, "", err.Error(), "")
		ctx.AbortWithStatusJSON(response.Status, response)
		return
	}

	gitRepoDto, err := h.gitRepoSvc.UpdateGitRepository(repoId, updateGitRepoDto)
	if err != nil {
		response := dto.NewResponseDto(http.StatusInternalServerError, "", err.Error(), "")
		ctx.AbortWithStatusJSON(response.Status, response)
		return
	}

	response := dto.NewResponseDto(http.StatusOK, "", "", gitRepoDto)
	ctx.JSON(response.Status, response)
}

func (h gitRepositoryHandler) HandleGetAllGitRepository(ctx *gin.Context) {
	gitRepoDtos, err := h.gitRepoSvc.GetAllGitRepository()
	if err != nil {
		response := dto.NewResponseDto(http.StatusInternalServerError, "", err.Error(), "")
		ctx.AbortWithStatusJSON(response.Status, response)
		return
	}

	response := dto.NewResponseDto(http.StatusOK, "", "", gitRepoDtos)
	ctx.JSON(response.Status, response)
}

func (h gitRepositoryHandler) HandleDeleteGitRepository(ctx *gin.Context) {
	repoId := ctx.Param("id")
	err := h.gitRepoSvc.DeleteGitRepository(repoId)
	if err != nil {
		response := dto.NewResponseDto(http.StatusInternalServerError, "", err.Error(), "")
		ctx.AbortWithStatusJSON(response.Status, response)
		return
	}

	response := dto.NewResponseDto(http.StatusOK, "GitRepository deleted successfully", "", "")
	ctx.JSON(response.Status, response)
}
