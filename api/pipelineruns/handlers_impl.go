package pipelineruns

import (
	"github.com/gin-gonic/gin"
)

type handler struct {
	svc Service
}

func NewHandler(svc Service) Handler {
	return &handler{svc: svc}
}

func (h *handler) HandleGetPipelineRuns(c *gin.Context) {
	response := h.svc.GetPipelineRuns()
	if response.Error != "" {
		c.AbortWithStatusJSON(response.Status, response)
	}
	c.JSON(response.Status, response)
}

func (h *handler) HandleUpdatePipelineRun(ctx *gin.Context) {
	response := h.svc.UpdatePipelineRun(ctx)
	if response.Error != "" {
		ctx.AbortWithStatusJSON(response.Status, response)
	}

	ctx.JSON(response.Status, response)
}
