package utils

import (
	"github.com/forhomme/app-base/errs"
	"user-management/app/user/domain"
)

func ParseMessage(err error) domain.ErrorMessage {
	return domain.ErrorMessage{
		Status:  errs.GetType(err),
		Message: err.Error(),
	}
}

func ParseResponse(status int, message string, data interface{}) domain.ResponseMessage {
	return domain.ResponseMessage{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
