package dto

type ResponseDto struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

func NewResponseDto(status int, message, error string, data interface{}) ResponseDto {
	return ResponseDto{
		Status:  status,
		Message: message,
		Data:    data,
		Error:   error,
	}
}
