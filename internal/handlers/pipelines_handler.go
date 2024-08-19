package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"spyrosmoux/api/internal/models"
	"spyrosmoux/api/internal/services"
)

type PipelinesHandler struct {
	pipelinesService services.PipelinesService
}

func NewPipelinesHandler(pipelinesService services.PipelinesService) *PipelinesHandler {
	return &PipelinesHandler{pipelinesService: pipelinesService}
}

func (pipelinesHandler *PipelinesHandler) CreatePipeline(c *gin.Context) {
	var pipeline models.Pipeline
	err := c.BindJSON(&pipeline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdPipeline, err := pipelinesHandler.pipelinesService.Create(&pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPipeline)
}

func (pipelinesHandler PipelinesHandler) FindPipelineById(c *gin.Context) {
	id := c.Param("id")
	pipeline, err := pipelinesHandler.pipelinesService.FindById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if pipeline == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pipeline not found"})
		return
	}

	c.JSON(http.StatusOK, pipeline)
}

func (pipelinesHandler *PipelinesHandler) FindAllPipelines(c *gin.Context) {
	pipelines, err := pipelinesHandler.pipelinesService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pipelines)
}

func (pipelinesHandler PipelinesHandler) DeletePipeline(c *gin.Context) {
	id := c.Param("id")
	err := pipelinesHandler.pipelinesService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "success")
}
