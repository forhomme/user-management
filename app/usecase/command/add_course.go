package command

import (
	"context"
	"fmt"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/mitchellh/mapstructure"
	"user-management/app/common/decorator"
	"user-management/app/domain/course"
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
	IsAssignment bool     `json:"is_assignment"`
	Ordering     int      `json:"ordering"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	FileUrl      []string `json:"file_url"`
}

type AddCourseHandler decorator.CommandHandler[*AddCourse]

type addCourseRepository struct {
	dbRepo course.CommandRepository
	logger logger.Logger
}

func NewAddCourseRepository(dbRepo course.CommandRepository, logger logger.Logger) decorator.CommandHandler[*AddCourse] {
	return decorator.ApplyCommandDecorators[*AddCourse](
		addCourseRepository{dbRepo: dbRepo, logger: logger},
		logger,
	)
}

func (a addCourseRepository) Handle(ctx context.Context, in *AddCourse) (err error) {
	var parentCourse *course.CoursePath
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
