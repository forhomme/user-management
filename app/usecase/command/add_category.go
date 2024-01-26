package command

import (
	"context"
	"github.com/forhomme/app-base/usecase/logger"
	"user-management/app/common/decorator"
	"user-management/app/domain/course"
)

type AddCategory struct {
	CategoryName string `json:"category_name"`
}

type AddCategoryHandler decorator.CommandHandler[*AddCategory]

type addCategoryRepository struct {
	dbRepo course.CommandRepository
	log    logger.Logger
}

func NewAddCategoryRepository(dbRepo course.CommandRepository, log logger.Logger) decorator.CommandHandler[*AddCategory] {
	return decorator.ApplyCommandDecorators[*AddCategory](
		addCategoryRepository{dbRepo: dbRepo, log: log},
		log)
}

func (a addCategoryRepository) Handle(ctx context.Context, in *AddCategory) error {
	err := a.dbRepo.AddCategory(ctx, in.CategoryName)
	if err != nil {
		a.log.Error(err)
		return err
	}
	return nil
}
