package pipelineruns

import (
	"errors"
	"strings"
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
