package _interface

import (
	"github.com/forhomme/app-base/usecase/controller"
	"user-management/app/domain/course"
)

var (
	content  = "content"
	filename = "filename"
)

func newUserFromContext(c controller.Context) (*course.User, error) {
	userId := c.GetAuthUser().UserId
	roleId := c.GetAuthUser().RoleId

	userRole, err := course.NewUserRoleFromRoleId(roleId)
	if err != nil {
		return nil, err
	}

	user := course.MustNewUser(userId, userRole)
	return &user, nil
}
