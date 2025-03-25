package dto

type LogEntryDto struct {
	RunId     string `json:"runId"`
	Timestamp string `json:"timestamp"`
	LogLevel  string `json:"logLevel"`
	Message   string `json:"message"`
}
