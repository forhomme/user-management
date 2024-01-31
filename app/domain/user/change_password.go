package user

import (
	"golang.org/x/crypto/bcrypt"
)

type ChangePassword struct {
	UserId      string `json:"user_id"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

func (c *ChangePassword) IsValid(user *User) bool {
	if c.OldPassword == c.NewPassword {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.OldPassword))
	if err != nil {
		return false
	}
	return true
}
