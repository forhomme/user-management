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

type UpdateCourseHandler decorator.CommandHandler[*AddCourse]

type updateCourseRepository struct {
	dbRepo adapters.CommandRepository
	logger logger.Logger
}

func NewUpdateCourseRepository(dbRepo adapters.CommandRepository, logger logger.Logger) decorator.CommandHandler[*AddCourse] {
	return decorator.ApplyCommandDecorators[*AddCourse](
		updateCourseRepository{dbRepo: dbRepo, logger: logger},
		logger,
	)
}

func (u updateCourseRepository) Handle(ctx context.Context, in *AddCourse) error {
	cm := &model.CoursePath{}
	err := mapstructure.Decode(in, cm)
	if err != nil {
		u.logger.Error(fmt.Errorf("error decode: %w", err))
		return err
	}

	err = u.dbRepo.ReplaceCourse(ctx, in.CourseId, cm)
	if err != nil {
		u.logger.Error(fmt.Errorf("error replace course: %w", err))
		return err
	}
	return nil
}
