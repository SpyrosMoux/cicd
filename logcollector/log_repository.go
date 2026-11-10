package logcollector

import (
	"gorm.io/gorm"
)

type LogRepository interface {
	Save(LogEntry) error
	FindAllByRunId(runId string) ([]LogEntry, error)
}

type logRepository struct {
	DbConn *gorm.DB
}

func NewLogRepository(dbConn *gorm.DB) LogRepository {
	return &logRepository{DbConn: dbConn}
}

func (logRepository logRepository) Save(logEntry LogEntry) error {
	result := logRepository.DbConn.Create(logEntry)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (logRepository logRepository) FindAllByRunId(runId string) ([]LogEntry, error) {
	var logEntries []LogEntry
	// Stored timestamps can repeat; UUIDs provide a stable tie-breaker, not emission order.
	result := logRepository.DbConn.Where("run_id = ?", runId).
		Order("timestamp ASC").Order("id ASC").Find(&logEntries)
	if result.Error != nil {
		return []LogEntry{}, result.Error
	}

	return logEntries, nil
}
