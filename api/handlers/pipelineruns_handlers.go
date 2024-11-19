package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/cicd/api/entities"
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

// TODO(@SpyrosMoux) improve logic, separate bussiness logic into service
func (h *pipelineRunsHandler) HandleUpdatePipelineRun(c *gin.Context) {
	runId := c.Param("id")

	var pipelineRun *entities.PipelineRun
	err := c.ShouldBindJSON(&pipelineRun)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	updatedPipelineRun, err := h.svc.UpdatePipelineRun(runId, pipelineRun)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, updatedPipelineRun)
}

// TODO(@SpyrosMoux) improve logic, separate bussiness logic into service
func (h *pipelineRunsHandler) HandleUpdatePipelineRunStatus(c *gin.Context) {
	runId := c.Param("id")

	var statusStr entities.StatusDto
	err := c.ShouldBindJSON(&statusStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	status, err := entities.ParseStatus(statusStr.Status)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	updatedPipelineRun, err := h.svc.UpdatePipelineRunStatus(runId, &status)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, updatedPipelineRun)
}
