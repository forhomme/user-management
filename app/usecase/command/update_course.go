package command

import (
	"context"
	"fmt"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/mitchellh/mapstructure"
	"user-management/app/common/decorator"
	"user-management/app/domain/course"
)

type ReplaceCourseHandler decorator.CommandHandler[*AddCourse]

type replaceCourseRepository struct {
	dbRepo course.CommandRepository
	logger logger.Logger
}

func NewReplaceCourseRepository(dbRepo course.CommandRepository, logger logger.Logger) decorator.CommandHandler[*AddCourse] {
	return decorator.ApplyCommandDecorators[*AddCourse](
		replaceCourseRepository{dbRepo: dbRepo, logger: logger},
		logger,
	)
}

func (u replaceCourseRepository) Handle(ctx context.Context, in *AddCourse) error {
	var cp course.CoursePath
	err := mapstructure.Decode(in, &cp)
	if err != nil {
		u.logger.Error(fmt.Errorf("error decode: %w", err))
		return err
	}

	return u.dbRepo.UpdateCourse(
		context.TODO(),
		in.CourseId,
		func(ctx context.Context, cm *course.CoursePath) (*course.CoursePath, error) {
			cm.Replace(cp)
			return cm, nil
		})
}
