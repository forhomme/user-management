package course

import (
	"errors"
	"fmt"
	"user-management/app/common/utils"
)

type UserRole struct {
	role string
}

func (u *UserRole) IsZero() bool {
	return *u == UserRole{}
}

func (u *UserRole) String() string {
	return u.role
}

var (
	Teacher = UserRole{"teacher"}
	Student = UserRole{"student"}
)

func NewUserRoleFromRoleId(roleId int) (UserRole, error) {
	switch roleId {
	case 1:
		return Teacher, nil
	case 2:
		return Student, nil
	}

	return UserRole{}, utils.NewIncorrectInputError(
		fmt.Sprintf("invalid '%d' role", roleId),
		"invalid-role",
	)
}

type User struct {
	userId string
	role   UserRole
}

func (u *User) UUID() string {
	return u.userId
}

func (u *User) Type() UserRole {
	return u.role
}

func (u *User) IsEmpty() bool {
	return *u == User{}
}

func NewUser(userUUID string, role UserRole) (User, error) {
	if userUUID == "" {
		return User{}, errors.New("missing user UUID")
	}
	if role.IsZero() {
		return User{}, errors.New("missing user type")
	}

	return User{userId: userUUID, role: role}, nil
}

func MustNewUser(userUUID string, role UserRole) User {
	u, err := NewUser(userUUID, role)
	if err != nil {
		panic(err)
	}

	return u
}

type ForbiddenToSeeCourseError struct {
	RequestingUserId string
}

func (f *ForbiddenToSeeCourseError) Error() string {
	return fmt.Sprintf(
		"user '%s' can't see course",
		f.RequestingUserId,
	)
}

func CanUserSeeCourse(user *User, course *CoursePath) error {
	if user.Type() == Teacher {
		return nil
	}

	if course.IsCourseVisible() {
		return nil
	}

	return &ForbiddenToSeeCourseError{user.UUID()}
}
