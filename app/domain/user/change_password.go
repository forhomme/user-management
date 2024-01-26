package user

import "user-management/app/common/utils"

type ChangePassword struct {
	UserId      string `json:"user_id"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

func (c *ChangePassword) IsValid(user *User) bool {
	if c.OldPassword == c.NewPassword {
		return false
	}
	hashOldPass, _ := utils.HashPassword(c.OldPassword)
	if user.Password != hashOldPass {
		return false
	}
	return true
}
