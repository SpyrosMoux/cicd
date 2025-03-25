package logcollector

import (
	"github.com/sirupsen/logrus"
	"github.com/spyrosmoux/cicd/common/dto"
)

type LogService interface {
	AddLog(dto.LogEntryDto) (LogEntry, error)
	GetLogsByRunId(runId string) ([]LogEntry, error)
}

type logService struct {
	Logger        *logrus.Logger
	LogRepository LogRepository
}

func NewLogService(logRepo LogRepository) LogService {
	return &logService{LogRepository: logRepo}
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

func logEntryDtoToLogEntry(logEntryDto dto.LogEntryDto) (LogEntry, error) {
	logLevel, err := fromString(logEntryDto.LogLevel)
	if err != nil {
		return LogEntry{}, err
	}
	return NewLogEntry(logEntryDto.RunId, logEntryDto.Timestamp, logEntryDto.Message, logLevel), nil
}
