package usecase

import (
	"user-management/app/courses/usecase/command"
	"user-management/app/courses/usecase/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AddCategory  command.AddCategoryHandler
	AddCourse    command.AddCourseHandler
	UpdateCourse command.UpdateCourseHandler
}

type Queries struct {
	GetAllCategories query.GetCategoryHandler
	GetCourses       query.GetCourseHandler
}
