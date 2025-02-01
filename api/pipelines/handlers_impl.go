package pipelines

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/cicd/common/dto"
)

type pipelineHandler struct {
	pipelineService PipelineService
}

func NewPipelineHandler(pipelineService PipelineService) PipelineHandler {
	return &pipelineHandler{pipelineService: pipelineService}
}

func (pipelineHandler pipelineHandler) HandleCreatePipeline(ctx *gin.Context) {
	var createPipelineDto CreatePipelineDto
	err := ctx.ShouldBindJSON(&createPipelineDto)
	if err != nil {
		response := dto.NewResponseDto(http.StatusBadRequest, "", err.Error(), "")
		ctx.AbortWithStatusJSON(response.Status, response)
		return
	}

	pipelineDto, err := pipelineHandler.pipelineService.CreatePipeline(createPipelineDto)
	if err != nil {
		response := dto.NewResponseDto(http.StatusInternalServerError, "", err.Error(), "")
		ctx.AbortWithStatusJSON(response.Status, response)
		return
	}

	response := dto.NewResponseDto(http.StatusCreated, "", "", pipelineDto)
	ctx.JSON(response.Status, response)
}

func (pipelineHandler pipelineHandler) HandleGetAllPipelines(ctx *gin.Context) {
	pipelines, err := pipelineHandler.pipelineService.GetAllPipelines()
	if err != nil {
		response := dto.NewResponseDto(http.StatusInternalServerError, "", err.Error(), "")
		ctx.AbortWithStatusJSON(response.Status, response)
		return
	}

	response := dto.NewResponseDto(http.StatusOK, "", "", pipelines)
	ctx.JSON(response.Status, response)
}

func (pipelineHandler pipelineHandler) HandleDeletePipeline(ctx *gin.Context) {
}
