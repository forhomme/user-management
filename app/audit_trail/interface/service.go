package _interface

import (
	"context"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/mitchellh/mapstructure"
	"user-management/app/audit_trail/usecase"
	"user-management/app/audit_trail/usecase/command"
	"user-management/config"
)

type AuditTrailService interface {
	InsertAuditTrail(c context.Context, in *AuditTrail) (err error)
}

type AuditTrailHandler struct {
	cfg    *config.Config
	logger logger.Logger
	app    usecase.Application
}

func NewAuditTrailService(cfg *config.Config, log logger.Logger, app usecase.Application) AuditTrailService {
	return AuditTrailHandler{cfg: cfg, logger: log, app: app}
}

func (a AuditTrailHandler) InsertAuditTrail(c context.Context, in *AuditTrail) (err error) {
	var addAuditTrail command.AddAuditTrail
	err = mapstructure.Decode(in, &addAuditTrail)
	if err != nil {
		a.logger.Error(err)
		return err
	}

	err = a.app.Commands.AddAuditTrail.Handle(c, &addAuditTrail)
	if err != nil {
		a.logger.Error(err)
		return err
	}

	return nil
}
