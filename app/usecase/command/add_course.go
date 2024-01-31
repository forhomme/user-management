package command

import (
	"context"
	"github.com/forhomme/app-base/infrastructure/telemetry"
	"github.com/forhomme/app-base/usecase/logger"
	"user-management/app/common/decorator"
	"user-management/app/domain/course"
)

type AddCourseHandler decorator.CommandHandler[*course.CoursePath]

type addCourseRepository struct {
	dbRepo course.CommandRepository
}

func NewAddCourseRepository(dbRepo course.CommandRepository, logger logger.Logger, tracer *telemetry.OtelSdk) decorator.CommandHandler[*course.CoursePath] {
	return decorator.ApplyCommandDecorators[*course.CoursePath](
		addCourseRepository{dbRepo: dbRepo},
		logger,
		tracer,
	)
}

func (a addCourseRepository) Handle(ctx context.Context, in *course.CoursePath) (err error) {
	err = a.dbRepo.AddCourse(ctx, in)
	if err != nil {
		return err
	}
	return nil
}
