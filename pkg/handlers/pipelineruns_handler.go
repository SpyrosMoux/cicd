package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/api/pkg/pipelineruns"
)

func HandleGetPipelineRuns(c *gin.Context) {
	client, err := pipelineruns.NewClient()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	runs, err := client.GetPipelineRuns()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, runs)
}

func UpdatePipelineRun(c *gin.Context) {
	runId := c.Param("name")

	client, err := pipelineruns.NewClient()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	var pipelineRun *pipelineruns.PipelineRun
	c.ShouldBindBodyWithJSON(&pipelineRun)

	updatedPipelineRun, err := client.UpdatePipelineRun(runId, pipelineRun)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, updatedPipelineRun)
}
