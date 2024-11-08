package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/cicd/api/pipelineruns"
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
	err := c.ShouldBindJSON(&pipelineRun)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	updatedPipelineRun, err := pipelineruns.UpdatePipelineRun(runId, pipelineRun)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, updatedPipelineRun)
}

func UpdatePipelineRunStatus(c *gin.Context) {
	runId := c.Param("id")

	var statusStr pipelineruns.StatusDto
	err := c.ShouldBindJSON(&statusStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	status, err := pipelineruns.ParseStatus(statusStr.Status)
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
