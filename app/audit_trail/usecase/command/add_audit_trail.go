package command

import (
	"context"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/mitchellh/mapstructure"
	"user-management/app/audit_trail/adapters"
	"user-management/app/courses/common/decorator"
)

type AddAuditTrail struct {
	UserId   string
	Menu     string
	Path     string
	Request  string
	Response string
}

type AddAuditTrailHandler decorator.CommandHandler[*AddAuditTrail]

type addAuditTrailRepository struct {
	dbRepo adapters.CommandRepository
	log    logger.Logger
}

func NewAddAuditTrailRepository(dbRepo adapters.CommandRepository, log logger.Logger) decorator.CommandHandler[*AddAuditTrail] {
	return decorator.ApplyCommandDecorators[*AddAuditTrail](
		addAuditTrailRepository{dbRepo: dbRepo, log: log},
		log)
}

func (a addAuditTrailRepository) Handle(ctx context.Context, in *AddAuditTrail) error {
	var auditTrail adapters.AuditTrailModel
	err := mapstructure.Decode(in, &auditTrail)
	if err != nil {
		a.log.Error(err)
		return err
	}
	err = a.dbRepo.AddAuditTrail(ctx, &auditTrail)
	if err != nil {
		a.log.Error(err)
		return err
	}
	return nil
}
