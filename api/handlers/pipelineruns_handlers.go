package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/cicd/api/services"
	"log"
	"net/http"
)

type PipelineRunsHandler interface {
	HandleGetPipelineRuns(c *gin.Context)
	HandleUpdatePipelineRun(c *gin.Context)
	HandleUpdatePipelineRunStatus(c *gin.Context)
}

type pipelineRunsHandler struct {
	svc services.PipelineRunsService
}

func NewPipelineRunsHandler(svc services.PipelineRunsService) PipelineRunsHandler {
	return &pipelineRunsHandler{svc: svc}
}

func (h *pipelineRunsHandler) HandleGetPipelineRuns(c *gin.Context) {
	runs, err := h.svc.GetPipelineRuns()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, runs)
}

func (h *pipelineRunsHandler) HandleUpdatePipelineRun(ctx *gin.Context) {
	updatedPipelineRun, err := h.svc.UpdatePipelineRun(ctx)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, updatedPipelineRun)
}

func (h *pipelineRunsHandler) HandleUpdatePipelineRunStatus(ctx *gin.Context) {
	updatedPipelineRun, err := h.svc.UpdatePipelineRunStatus(ctx)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, updatedPipelineRun)
}
