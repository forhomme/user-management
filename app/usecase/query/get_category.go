package query

import (
	"context"
	"fmt"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/mitchellh/mapstructure"
	"user-management/app/common/decorator"
	"user-management/app/domain/course"
)

type AllCategory struct {
	Category []*Category `json:"category"`
}

type Category struct {
	CategoryId   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type GetCategoryHandler decorator.QueryHandler[*Category, *AllCategory]

type getCategoryRepository struct {
	dbRepo course.QueryRepository
	log    logger.Logger
}

func NewGetCategoryRepository(dbRepo course.QueryRepository, log logger.Logger) decorator.QueryHandler[*Category, *AllCategory] {
	return decorator.ApplyQueryDecorators[*Category, *AllCategory](
		getCategoryRepository{dbRepo: dbRepo, log: log},
		log)
}

func (g getCategoryRepository) Handle(ctx context.Context, _ *Category) (*AllCategory, error) {
	data, err := g.dbRepo.GetCategories(ctx)
	if err != nil {
		g.log.Error(fmt.Errorf("error get category from database: %w", err))
		return nil, err
	}

	out := make([]*Category, 0)
	err = mapstructure.Decode(data, &out)
	if err != nil {
		g.log.Error(fmt.Errorf("error decode category: %w", err))
		return nil, err
	}

	return &AllCategory{Category: out}, nil
}
