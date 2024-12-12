package pipelineruns

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type handler struct {
	svc Service
}

func NewHandler(svc Service) Handler {
	return &handler{svc: svc}
}

func (h *handler) HandleGetPipelineRuns(c *gin.Context) {
	runs, err := h.svc.GetPipelineRuns()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, runs)
}

func (h *handler) HandleUpdatePipelineRun(ctx *gin.Context) {
	updatedPipelineRun, err := h.svc.UpdatePipelineRun(ctx)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, updatedPipelineRun)
}
