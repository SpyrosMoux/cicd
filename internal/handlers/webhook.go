package handlers

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func HandleWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading body"})
		log.Panicf("Error reading body: %v", err)
		return
	}

	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty body"})
		log.Panicf("Error reading body: %v", err)
		return
	}

	var result map[string]interface{}
	err = c.ShouldBindBodyWithJSON(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Panicf("Error unmarshalling body: %v", err)
	}

	log.Println(body)

	//git.PlainClone(result.Repository)

	// Publish the raw YAML body as a job
	// file, _ := os.ReadFile("sample-pipeline.yaml") // TODO use body instead
	// queue.PublishJob(string(file))

	c.Status(http.StatusAccepted)
}
