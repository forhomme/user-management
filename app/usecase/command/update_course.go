package command

import (
	"context"
	"github.com/forhomme/app-base/infrastructure/baselogger"
	"github.com/forhomme/app-base/infrastructure/telemetry"
	"user-management/app/common/decorator"
	"user-management/app/domain/course"
)

type ReplaceCourseHandler decorator.CommandHandler[*course.CoursePath]

type replaceCourseRepository struct {
	dbRepo course.CommandRepository
}

func NewReplaceCourseRepository(dbRepo course.CommandRepository, logger *baselogger.Logger, tracer *telemetry.OtelSdk) decorator.CommandHandler[*course.CoursePath] {
	return decorator.ApplyCommandDecorators[*course.CoursePath](
		replaceCourseRepository{dbRepo: dbRepo},
		logger,
		tracer,
	)
}

func (u replaceCourseRepository) Handle(ctx context.Context, in *course.CoursePath) error {
	return u.dbRepo.UpdateCourse(
		ctx,
		in.CourseId,
		func(ctx context.Context, cm *course.CoursePath) (*course.CoursePath, error) {
			cm.Replace(in)
			return cm, nil
		})
}
