package logcollector

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/cicd/common/dto"
)

type Handler interface {
	HandleHealth(ctx *gin.Context)
	HandleGetLogsByRunId(ctx *gin.Context)
}

type handler struct {
	logSvc LogService
}

func NewHandler(logSvc LogService) Handler {
	return &handler{
		logSvc: logSvc,
	}
}

func (handler handler) HandleHealth(ctx *gin.Context) {
	response := dto.NewResponseDto(http.StatusOK, "I'm Alive!", "", "")
	ctx.JSON(response.Status, response)
}

func (handler handler) HandleGetLogsByRunId(ctx *gin.Context) {
	runId := ctx.Param("runId")
	logs, err := handler.logSvc.GetLogsByRunId(runId)
	if err != nil {
		response := dto.NewResponseDto(http.StatusInternalServerError, "", err.Error(), nil)
		ctx.AbortWithStatusJSON(response.Status, response)
		return
	}

	response := dto.NewResponseDto(http.StatusOK, "", "", logs)
	ctx.JSON(http.StatusOK, response)
}
