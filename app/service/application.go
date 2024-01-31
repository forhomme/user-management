package service

import (
	"github.com/forhomme/app-base/infrastructure/telemetry"
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

func NewApplication(cfg *config.Config, log logger.Logger, client *mongo.Client, sqlHandler database.SqlHandler, storage storage.Storage, tracer *telemetry.OtelSdk) usecase.Application {
	mongoDB := client.Database(cfg.Database)
	mongoRepo := adapters.NewCourseMongoRepository(cfg, log, mongoDB, tracer)
	mysqlRepo := adapters.NewCourseMysqlRepository(cfg, log, sqlHandler, tracer)

	return usecase.Application{
		Commands: usecase.Commands{
			ChangePassword: command.NewChangePasswordRepository(mysqlRepo, log, tracer),
			AddCategory:    command.NewAddCategoryRepository(mysqlRepo, log, tracer),
			AddCourse:      command.NewAddCourseRepository(mongoRepo, log, tracer),
			ReplaceCourse:  command.NewReplaceCourseRepository(mongoRepo, log, tracer),
		},
		Queries: usecase.Queries{
			SignUp:           query.NewSignUpRepository(mysqlRepo, log, cfg, tracer),
			Login:            query.NewLoginRepository(mysqlRepo, log, cfg, tracer),
			RefreshToken:     query.NewRefreshTokenRepository(mysqlRepo, log, cfg, tracer),
			GetAllCategories: query.NewGetCategoryRepository(mysqlRepo, log, tracer),
			GetCourses:       query.NewGetCourseRepository(mongoRepo, log, tracer),
			UploadFile:       query.NewUploadFileHandler(cfg, storage, log, tracer),
		},
	}
}
