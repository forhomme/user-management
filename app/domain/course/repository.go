package course

import (
	"context"
)

type CommandRepository interface {
	AddCategory(ctx context.Context, categoryName string) error
	AddCourse(ctx context.Context, cr *CoursePath) error
	UpdateCourse(ctx context.Context, id string, updateFn func(ctx context.Context, cm *CoursePath) (*CoursePath, error)) error
}

type QueryRepository interface {
	GetCourses(ctx context.Context, in *FilterCourse) ([]*CoursePath, error)
	GetCategories(ctx context.Context) ([]*Category, error)
}
