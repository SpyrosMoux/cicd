package pipelineruns

import (
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	PENDING Status = iota
	RUNNING
	CANCELED
	FAILED
	COMPLETED
)

func (s Status) String() string {
	switch s {
	case PENDING:
		return "Pending"
	case RUNNING:
		return "Running"
	case CANCELED:
		return "Canceled"
	case FAILED:
		return "Failed"
	case COMPLETED:
		return "Completed"
	default:
		return "Unknown"
	}
}

type PipelineRun struct {
	Id            uuid.UUID `json:"id" gorm:"primaryKey"`
	Status        string    `json:"status"`
	Repository    string    `json:"repository"`
	Branch        string    `json:"branch"`
	TimeTriggered int64     `json:"time_triggered"`
	TimeStarted   int64     `json:"time_started"`
	TimeEnded     int64     `json:"time_ended"`
}

func NewPipelineRun(repository, branch string) *PipelineRun {
	return &PipelineRun{
		Id:            uuid.New(),
		Status:        PENDING.String(),
		Repository:    repository,
		Branch:        branch,
		TimeTriggered: time.Now().Unix(),
		TimeStarted:   0, // init as 0, will be updated by the runner once started
		TimeEnded:     0, // init as 0, will be updated by the runner once finished
	}
}
