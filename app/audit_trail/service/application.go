package service

import (
	"github.com/forhomme/app-base/usecase/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"user-management/app/audit_trail/adapters"
	"user-management/app/audit_trail/usecase"
	"user-management/app/audit_trail/usecase/command"
	"user-management/config"
)

func NewAuditApplication(cfg *config.Config, log logger.Logger, client *mongo.Client) usecase.Application {
	mongoDB := client.Database(cfg.Database)
	dbRepo := adapters.NewCourseMongoRepository(cfg, log, mongoDB)

	return usecase.Application{
		Commands: usecase.Commands{
			AddAuditTrail: command.NewAddAuditTrailRepository(dbRepo, log),
		},
	}
}
