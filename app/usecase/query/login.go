package query

import (
	"context"
	"errors"
	"github.com/forhomme/app-base/infrastructure/baselogger"
	"github.com/forhomme/app-base/infrastructure/telemetry"
	"user-management/app/common/decorator"
	"user-management/app/domain/user"
	"user-management/config"
)

type LoginHandler decorator.QueryHandler[*user.Login, *user.Token]

type loginRepository struct {
	dbRepo user.QueryRepository
	cfg    *config.Config
}

func NewLoginRepository(dbRepo user.QueryRepository, logger *baselogger.Logger, cfg *config.Config, tracer *telemetry.OtelSdk) decorator.QueryHandler[*user.Login, *user.Token] {
	return decorator.ApplyQueryDecorators[*user.Login, *user.Token](
		loginRepository{dbRepo: dbRepo, cfg: cfg},
		logger,
		tracer,
	)
}

func (l loginRepository) Handle(ctx context.Context, in *user.Login) (token *user.Token, err error) {
	userData, err := l.dbRepo.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return nil, err
	}
	if !userData.IsExist() {
		return nil, errors.New("user not found")
	}

	err = in.CheckPassword(userData.Password)
	if err != nil {
		return nil, err
	}

	token, err = userData.GenerateToken(l.cfg)
	if err != nil {
		return nil, err
	}
	return token, nil
}
