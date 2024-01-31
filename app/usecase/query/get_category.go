package query

import (
	"context"
	"github.com/forhomme/app-base/infrastructure/telemetry"
	"github.com/forhomme/app-base/usecase/logger"
	"user-management/app/common/decorator"
	"user-management/app/domain/course"
)

type AllCategory struct {
	Category []*course.Category `json:"category"`
}

type GetCategoryHandler decorator.QueryHandler[*course.Category, *AllCategory]

type getCategoryRepository struct {
	dbRepo course.QueryRepository
}

func NewGetCategoryRepository(dbRepo course.QueryRepository, log logger.Logger, tracer *telemetry.OtelSdk) decorator.QueryHandler[*course.Category, *AllCategory] {
	return decorator.ApplyQueryDecorators[*course.Category, *AllCategory](
		getCategoryRepository{dbRepo: dbRepo},
		log,
		tracer)
}

func (g getCategoryRepository) Handle(ctx context.Context, _ *course.Category) (out *AllCategory, err error) {
	categories, err := g.dbRepo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	return &AllCategory{Category: categories}, nil
}
