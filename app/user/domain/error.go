package domain

import "github.com/forhomme/app-base/errs"

type ErrorMessage struct {
	Status  errs.ErrorType
	Message string
}

type ResponseMessage struct {
	Status  int
	Message string
	Data    interface{} `json:"Data,omitempty"`
}
