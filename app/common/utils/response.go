package utils

import (
	"github.com/forhomme/app-base/errs"
)

type ErrorMessage struct {
	Status  errs.ErrorType `json:"status"`
	Message string         `json:"message"`
}

func ParseMessage(err error) ErrorMessage {
	return ErrorMessage{
		Status:  errs.GetType(err),
		Message: err.Error(),
	}
}

type ResponseMessage struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ParseResponse(status int, message string, data interface{}) ResponseMessage {
	return ResponseMessage{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
