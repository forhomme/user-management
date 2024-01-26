package service

import (
	"github.com/forhomme/app-base/usecase/database"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/forhomme/app-base/usecase/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"user-management/app/adapters"
	"user-management/app/usecase"
	"user-management/app/usecase/command"
	"user-management/app/usecase/query"
	"user-management/config"
)

func NewApplication(cfg *config.Config, log logger.Logger, client *mongo.Client, sqlHandler database.SqlHandler, storage storage.Storage) usecase.Application {
	mongoDB := client.Database(cfg.Database)
	mongoRepo := adapters.NewCourseMongoRepository(cfg, log, mongoDB)
	mysqlRepo := adapters.NewCourseMysqlRepository(cfg, log, sqlHandler)

	return usecase.Application{
		Commands: usecase.Commands{
			ChangePassword: command.NewChangePasswordRepository(mysqlRepo, log),
			AddCategory:    command.NewAddCategoryRepository(mysqlRepo, log),
			AddCourse:      command.NewAddCourseRepository(mongoRepo, log),
			ReplaceCourse:  command.NewReplaceCourseRepository(mongoRepo, log),
		},
		Queries: usecase.Queries{
			SignUp:           query.NewSignUpRepository(mysqlRepo, log, cfg),
			Login:            query.NewLoginRepository(mysqlRepo, log, cfg),
			GetAllCategories: query.NewGetCategoryRepository(mysqlRepo, log),
			GetCourses:       query.NewGetCourseRepository(mongoRepo, log),
			UploadFile:       query.NewUploadFileHandler(cfg, storage, log),
		},
	}
}
