package service

import (
	"github.com/forhomme/app-base/usecase/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"user-management/app/courses/adapters"
	"user-management/app/courses/usecase"
	"user-management/app/courses/usecase/command"
	"user-management/app/courses/usecase/query"
	"user-management/config"
)

func NewApplication(cfg *config.Config, log logger.Logger, client *mongo.Client) usecase.Application {
	mongoDB := client.Database(cfg.Database)
	dbRepo := adapters.NewCourseMongoRepository(cfg, log, mongoDB)

	return usecase.Application{
		Commands: usecase.Commands{
			AddCategory:  command.NewAddCategoryRepository(dbRepo, log),
			AddCourse:    command.NewAddCourseRepository(dbRepo, log),
			UpdateCourse: command.NewUpdateCourseRepository(dbRepo, log),
		},
		Queries: usecase.Queries{
			GetAllCategories: query.NewGetCategoryRepository(dbRepo, log),
			GetCourses:       query.NewGetCourseRepository(dbRepo, log),
		},
	}
}
