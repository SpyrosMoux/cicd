package handlers

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"spyrosmoux/api/internal/queue"
)

func HandleWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading body"})
		return
	}

	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty body"})
		return
	}

	// Publish the raw YAML body as a job
	queue.PublishJob(string(body))

	c.Status(http.StatusAccepted)
}
