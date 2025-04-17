package logcollector

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spyrosmoux/cicd/common/dto"
)

type Handler interface {
	HandleHealth(ctx *gin.Context)
	HandleGetLogsByRunId(ctx *gin.Context)
	HandleStreamLogsByRunId(ctx *gin.Context)
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

func (handler handler) HandleStreamLogsByRunId(ctx *gin.Context) {
	runId := ctx.Param("runId")
	if runId == "" {
		response := dto.NewResponseDto(http.StatusBadRequest, "", "Missing runId", nil)
		ctx.AbortWithStatusJSON(response.Status, response)
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		response := dto.NewResponseDto(http.StatusInternalServerError, "", "Failed to upgrade connection", nil)
		ctx.AbortWithStatusJSON(response.Status, response)
		return
	}

	handler.logSvc.HandleWebSocket(runId, conn)

	response := dto.NewResponseDto(http.StatusOK, "", "", nil)
	ctx.JSON(http.StatusOK, response)
}
