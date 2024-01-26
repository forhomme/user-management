package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

var WrongUsernameOrPasswordError = errors.New("wrong username or password")

func (l *Login) CheckPassword(hashPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(l.Password))
	if err != nil {
		return WrongUsernameOrPasswordError
	}
	return nil
}
