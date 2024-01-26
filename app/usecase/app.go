package usecase

import (
	"user-management/app/usecase/command"
	"user-management/app/usecase/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	ChangePassword command.ChangePasswordHandler
	AddCategory    command.AddCategoryHandler
	AddCourse      command.AddCourseHandler
	ReplaceCourse  command.ReplaceCourseHandler
}

type Queries struct {
	SignUp           query.SignUpHandler
	Login            query.LoginHandler
	RefreshToken     query.RefreshTokenHandler
	GetAllCategories query.GetCategoryHandler
	GetCourses       query.GetCourseHandler
	UploadFile       query.UploadFileHandler
}
