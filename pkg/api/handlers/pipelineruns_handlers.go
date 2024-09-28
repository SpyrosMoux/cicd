package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/api/pkg/business/pipelineruns"
	"log"
	"net/http"
)

func HandleGetPipelineRuns(c *gin.Context) {
	runs, err := pipelineruns.GetPipelineRuns()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, runs)
}

func UpdatePipelineRun(c *gin.Context) {
	runId := c.Param("id")

	var pipelineRun *pipelineruns.PipelineRun
	c.ShouldBindBodyWithJSON(&pipelineRun)

	updatedPipelineRun, err := pipelineruns.UpdatePipelineRun(runId, pipelineRun)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, updatedPipelineRun)
}

func UpdatePipelineRunStatus(c *gin.Context) {
	runId := c.Param("id")
	statusStr := c.Param("status")

	status, err := pipelineruns.ParseStatus(statusStr)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	updatedPipelineRun, err := pipelineruns.UpdatePipelineRunStatus(runId, status)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, updatedPipelineRun)
}
