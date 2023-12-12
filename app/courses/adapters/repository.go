package adapters

import (
	"context"
	"user-management/app/courses/domain/model"
)

type CommandRepository interface {
	AddCategory(ctx context.Context, categoryName string) error
	AddCourse(ctx context.Context, cr *model.CoursePath) error
	UpdateCourse(ctx context.Context, id string, updateFn func(ctx context.Context, cm *model.CoursePath) (*model.CoursePath, error)) error
	ReplaceCourse(ctx context.Context, id string, in *model.CoursePath) error
}

type QueryRepository interface {
	GetCourses(ctx context.Context, in *FilterModel) ([]*model.CoursePath, error)
	GetCategories(ctx context.Context) ([]*CategoryModel, error)
}
