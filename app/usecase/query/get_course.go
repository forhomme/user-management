package query

import (
	"context"
	"github.com/forhomme/app-base/infrastructure/telemetry"
	"github.com/forhomme/app-base/usecase/logger"
	"user-management/app/common/decorator"
	"user-management/app/domain/course"
)

type AllCourse struct {
	Courses []*course.CoursePath `json:"Courses"`
}

type GetCourseHandler decorator.QueryHandler[*course.FilterCourse, *AllCourse]

type getCourseRepository struct {
	dbRepo course.QueryRepository
}

func NewGetCourseRepository(dbRepo course.QueryRepository, logger logger.Logger, tracer *telemetry.OtelSdk) decorator.QueryHandler[*course.FilterCourse, *AllCourse] {
	return decorator.ApplyQueryDecorators[*course.FilterCourse, *AllCourse](
		getCourseRepository{dbRepo: dbRepo},
		logger,
		tracer,
	)
}

func (g getCourseRepository) Handle(ctx context.Context, in *course.FilterCourse) (out *AllCourse, err error) {
	data := make([]*course.CoursePath, 0)
	courses, err := g.dbRepo.GetCourses(ctx, in)
	if err != nil {
		return nil, err
	}
	for _, c := range courses {
		if err = course.CanUserSeeCourse(in.User, c); err != nil {
			continue
		}
		c.List()
		data = append(data, c)
	}
	return &AllCourse{Courses: data}, nil
}
