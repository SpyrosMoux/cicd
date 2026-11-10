package logcollector

import (
	"fmt"

	"github.com/google/uuid"
)

type LogLevel int

const (
	INFO LogLevel = 0
	WARN
	ERROR
	DEBUG
)

type LogEntry struct {
	Id        string   `json:"id" gorm:"primary_key"`
	RunId     string   `json:"runid"`
	Timestamp string   `json:"timestamp"`
	LogLevel  LogLevel `json:"logLevel"`
	Message   string   `json:"message"`
}

func fromString(logLevel string) (LogLevel, error) {
	switch logLevel {
	case "INFO":
		return INFO, nil
	case "WARN":
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	case "DEBUG":
		return DEBUG, nil
	default:
		return 0, fmt.Errorf("unknown log level=%s", logLevel)
	}
}

func NewLogEntry(runId, timestamp, message string, logLevel LogLevel) LogEntry {
	return LogEntry{
		Id:        uuid.NewString(),
		RunId:     runId,
		Timestamp: timestamp,
		LogLevel:  logLevel,
		Message:   message,
	}
}
