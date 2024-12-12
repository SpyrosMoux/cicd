package dto

// UpdatePipelineRunDto used to update a pipeline run when it
// starts and when it finishes
type UpdatePipelineRunDto struct {
	Status      string `json:"status"`
	TimeStarted int64  `json:"time_started,omitempty"`
	TimeEnded   int64  `json:"time_ended,omitempty"`
}
