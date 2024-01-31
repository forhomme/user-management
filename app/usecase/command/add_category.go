package command

import (
	"context"
	"github.com/forhomme/app-base/infrastructure/telemetry"
	"github.com/forhomme/app-base/usecase/logger"
	"user-management/app/common/decorator"
	"user-management/app/domain/course"
)

type AddCategoryHandler decorator.CommandHandler[*course.Category]

type addCategoryRepository struct {
	dbRepo course.CommandRepository
}

func NewAddCategoryRepository(dbRepo course.CommandRepository, log logger.Logger, tracer *telemetry.OtelSdk) decorator.CommandHandler[*course.Category] {
	return decorator.ApplyCommandDecorators[*course.Category](
		addCategoryRepository{dbRepo: dbRepo},
		log,
		tracer)
}

func (a addCategoryRepository) Handle(ctx context.Context, in *course.Category) (err error) {
	err = a.dbRepo.AddCategory(ctx, in.CategoryName)
	if err != nil {
		return err
	}
	return nil
}
