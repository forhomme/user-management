package utils

import (
	"github.com/forhomme/app-base/errs"
)

type ErrorMessage struct {
	status  errs.ErrorType `json:"status"`
	message string         `json:"message"`
}

func ParseMessage(err error) ErrorMessage {
	return ErrorMessage{
		status:  errs.GetType(err),
		message: err.Error(),
	}
}

type ResponseMessage struct {
	status  int         `json:"status"`
	message string      `json:"message"`
	data    interface{} `json:"data"`
}

func ParseResponse(status int, message string, data interface{}) ResponseMessage {
	return ResponseMessage{
		status:  status,
		message: message,
		data:    data,
	}
}
