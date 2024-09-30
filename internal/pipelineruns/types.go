package pipelineruns

import (
	"github.com/spyrosmoux/api/pkg/pipelineruns"
	"time"

	"github.com/google/uuid"
)

type PipelineRun struct {
	Id            string `json:"id" gorm:"primary_key"`
	Status        string `json:"status"`
	Repository    string `json:"repository"`
	Branch        string `json:"branch"`
	TimeTriggered int64  `json:"time_triggered"`
	TimeStarted   int64  `json:"time_started"`
	TimeEnded     int64  `json:"time_ended"`
}

func NewPipelineRun(repository, branch string) *PipelineRun {
	return &PipelineRun{
		Id:            uuid.New().String(),
		Status:        pipelineruns.PENDING.String(),
		Repository:    repository,
		Branch:        branch,
		TimeTriggered: time.Now().Unix(),
		TimeStarted:   0, // init as 0, will be updated by the runner once started
		TimeEnded:     0, // init as 0, will be updated by the runner once finished
	}
}
