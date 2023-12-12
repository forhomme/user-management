package query

import (
	"context"
	"fmt"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/mitchellh/mapstructure"
	"user-management/app/courses/adapters"
	"user-management/app/courses/common/decorator"
	"user-management/app/courses/usecase/command"
)

type GetCourses struct {
	Id         string
	CategoryId int    `json:"CategoryId"`
	Filter     string `json:"Filter"`
	PerPage    int    `json:"PerPage"`
	Page       int    `json:"Page"`
}

type AllCourse struct {
	Courses []*command.AddCourse `json:"Courses"`
}

type GetCourseHandler decorator.QueryHandler[*GetCourses, *AllCourse]

type getCourseRepository struct {
	dbRepo adapters.QueryRepository
	logger logger.Logger
}

func NewGetCourseRepository(dbRepo adapters.QueryRepository, logger logger.Logger) decorator.QueryHandler[*GetCourses, *AllCourse] {
	return decorator.ApplyQueryDecorators[*GetCourses, *AllCourse](
		getCourseRepository{dbRepo: dbRepo, logger: logger},
		logger,
	)
}

func (g getCourseRepository) Handle(ctx context.Context, in *GetCourses) (*AllCourse, error) {
	out := make([]*command.AddCourse, 0)
	courses, err := g.dbRepo.GetCourses(ctx, &adapters.FilterModel{
		ID:         in.Id,
		Filter:     in.Filter,
		CategoryId: in.CategoryId,
		Page:       int64(in.Page),
		PerPage:    int64(in.PerPage),
	})
	if err != nil {
		g.logger.Error(fmt.Errorf("error get all data course %w", err))
		return nil, err
	}
	for _, course := range courses {
		data := &command.AddCourse{}
		err = mapstructure.Decode(course, data)
		if err != nil {
			g.logger.Error(fmt.Errorf("error decode: %w", err))
			continue
		}
		out = append(out, data)
	}
	return &AllCourse{Courses: out}, nil
}
