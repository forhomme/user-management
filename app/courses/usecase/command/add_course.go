package command

import (
	"context"
	"fmt"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/mitchellh/mapstructure"
	"user-management/app/courses/adapters"
	"user-management/app/courses/common/decorator"
	"user-management/app/courses/domain/model"
)

type AddCourse struct {
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

type AddCourseHandler decorator.CommandHandler[*AddCourse]

type addCourseRepository struct {
	dbRepo adapters.CommandRepository
	logger logger.Logger
}

func NewAddCourseRepository(dbRepo adapters.CommandRepository, logger logger.Logger) decorator.CommandHandler[*AddCourse] {
	return decorator.ApplyCommandDecorators[*AddCourse](
		addCourseRepository{dbRepo: dbRepo, logger: logger},
		logger,
	)
}

func (a addCourseRepository) Handle(ctx context.Context, in *AddCourse) (err error) {
	var parentCourse *model.CoursePath
	err = mapstructure.Decode(in, &parentCourse)
	if err != nil {
		a.logger.Error(fmt.Errorf("error decode input: %w", err))
		return err
	}
	err = a.dbRepo.AddCourse(ctx, parentCourse)
	if err != nil {
		a.logger.Error(fmt.Errorf("error add course: %w", err))
		return err
	}
	return nil
}
