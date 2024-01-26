package _interface

import (
	"github.com/forhomme/app-base/usecase/controller"
	"user-management/app/domain/course"
)

var (
	content  = "content"
	filename = "filename"
)

type Category struct {
	CategoryId   int    `json:"category_id"`
	CategoryName string `json:"category_name" validate:"required"`
}

type AllCategory struct {
	Category []*Category `json:"category"`
}

type GetCourses struct {
	Id         string
	CategoryId int    `json:"CategoryId"`
	Filter     string `json:"Filter"`
	PerPage    int    `json:"PerPage"`
	Page       int    `json:"Page"`
}

type AllCourse struct {
	Courses []*Course `json:"Courses"`
}

type Course struct {
	CourseId    string
	CategoryId  int
	Title       string
	Description string
	Tags        []string
	SubCourses  []*SubCourse
}

type SubCourse struct {
	Title       string
	Description string
	Contents    []*Content
}

type Content struct {
	IsAssignment bool   `json:"is_assignment"`
	Ordering     int    `json:"ordering"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Content      string `json:"content"`
}

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
