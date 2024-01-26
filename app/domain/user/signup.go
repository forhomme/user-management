package user

import (
	"errors"
	"user-management/app/common/utils"
)

type SignUp struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	RoleId   int    `json:"role_id" validate:"required"`
}

var errorGeneratePassword = errors.New("cannot hashing password")

func (s *SignUp) HashPassword() error {
	hashPass, err := utils.HashPassword(s.Password)
	if err != nil {
		return errorGeneratePassword
	}
	s.Password = hashPass
	return nil
}
