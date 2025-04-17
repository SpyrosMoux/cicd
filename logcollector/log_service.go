package logcollector

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spyrosmoux/cicd/common/dto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[string][]*websocket.Conn)
var mutex = &sync.Mutex{}

type LogService interface {
	AddLog(dto.LogEntryDto) (LogEntry, error)
	GetLogsByRunId(runId string) ([]LogEntry, error)
	HandleWebSocket(runId string, conn *websocket.Conn)
	BroadcastLog(runId, message string)
}

type logService struct {
	Logger        *logrus.Logger
	LogRepository LogRepository
}

func NewLogService(logRepo LogRepository) LogService {
	return &logService{
		LogRepository: logRepo,
	}
}

func (logSvc logService) AddLog(logEntryDto dto.LogEntryDto) (LogEntry, error) {
	logEntry, err := logEntryDtoToLogEntry(logEntryDto)
	if err != nil {
		logSvc.Logger.Error(err)
		return LogEntry{}, err
	}

	err = logSvc.LogRepository.Save(logEntry)
	if err != nil {
		logSvc.Logger.WithFields(logrus.Fields{
			"logId": logEntry.Id,
		}).WithError(err).Error("failed to save log")
		return LogEntry{}, err
	}

	// broadcast log to active websockets
	logSvc.BroadcastLog(logEntry.RunId, logEntry.Message)

	return logEntry, nil
}

func (logSvc logService) GetLogsByRunId(runId string) ([]LogEntry, error) {
	logEntries, err := logSvc.LogRepository.FindAllByRunId(runId)
	if err != nil {
		logSvc.Logger.WithFields(logrus.Fields{
			"runId": runId,
			"err":   err,
		}).Error("failed to find logs by run id")
		return []LogEntry{}, err
	}

	return logEntries, nil
}

func (logSvc logService) HandleWebSocket(runId string, conn *websocket.Conn) {
	mutex.Lock()
	clients[runId] = append(clients[runId], conn)
	mutex.Unlock()

	defer func() {
		conn.Close()
		mutex.Lock()
		for index, tmpConn := range clients[runId] {
			if tmpConn == conn {
				clients[runId] = append(clients[runId][:index], clients[runId][index+1:]...)
				break
			}
		}
		mutex.Unlock()
	}()

	// keeps the connection open
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (logSvc logService) BroadcastLog(runId, message string) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, conn := range clients[runId] {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			conn.Close()
		}
	}
}

func logEntryDtoToLogEntry(logEntryDto dto.LogEntryDto) (LogEntry, error) {
	logLevel, err := fromString(logEntryDto.LogLevel)
	if err != nil {
		return LogEntry{}, err
	}
	return NewLogEntry(logEntryDto.RunId, logEntryDto.Timestamp, logEntryDto.Message, logLevel), nil
}
