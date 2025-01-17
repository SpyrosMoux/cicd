package pipelineruns

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PipelineRun struct {
	Id            string `json:"id" gorm:"primary_key"`
	Status        string `json:"status"`
	Repository    string `json:"repository"`
	Branch        string `json:"branch"`
	TriggerType   string `json:"trigger_type"`
	TriggeredBy   string `json:"triggered_by"`
	TimeTriggered int64  `json:"time_triggered"`
	TimeStarted   int64  `json:"time_started"`
	TimeEnded     int64  `json:"time_ended"`
}

type TriggerType int

const (
	PUSH TriggerType = iota
	PR
)

func (t TriggerType) String() string {
	switch t {
	case PUSH:
		return "Push"
	case PR:
		return "Pull Request"
	default:
		return "Unknown"
	}
}

func NewPipelineRun(repository, branch, triggeredBy string, triggerType TriggerType) *PipelineRun {
	return &PipelineRun{
		Id:            uuid.New().String(),
		Status:        PENDING.String(),
		Repository:    repository,
		Branch:        branch,
		TriggerType:   triggerType.String(),
		TriggeredBy:   triggeredBy,
		TimeTriggered: time.Now().Unix(),
		TimeStarted:   0, // init as 0, will be updated by the runner once started
		TimeEnded:     0, // init as 0, will be updated by the runner once finished
	}
}

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

func ParseStatus(statusStr string) (Status, error) {
	switch strings.ToLower(statusStr) {
	case "pending":
		return PENDING, nil
	case "running":
		return RUNNING, nil
	case "canceled":
		return CANCELED, nil
	case "failed":
		return FAILED, nil
	case "completed":
		return COMPLETED, nil
	default:
		return 0, errors.New("invalid status")
	}
}

type StatusDto struct {
	Status string `json:"status"`
}
