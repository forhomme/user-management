package query

import (
	"context"
	"github.com/forhomme/app-base/infrastructure/baselogger"
	"github.com/forhomme/app-base/infrastructure/telemetry"
	"github.com/google/uuid"
	"user-management/app/common/decorator"
	"user-management/app/domain/user"
	"user-management/config"
)

type SignUpHandler decorator.QueryHandler[*user.SignUp, *user.Token]

type signUpRepository struct {
	dbRepo user.CommandRepository
	cfg    *config.Config
}

func NewSignUpRepository(dbRepo user.CommandRepository, logger *baselogger.Logger, cfg *config.Config, tracer *telemetry.OtelSdk) decorator.QueryHandler[*user.SignUp, *user.Token] {
	return decorator.ApplyQueryDecorators[*user.SignUp, *user.Token](
		signUpRepository{dbRepo: dbRepo, cfg: cfg},
		logger,
		tracer,
	)
}

func (s signUpRepository) Handle(ctx context.Context, in *user.SignUp) (token *user.Token, err error) {
	err = in.HashPassword()
	if err != nil {
		return nil, err
	}
	userData := &user.User{
		RoleId:   in.RoleId,
		UserId:   uuid.New().String(),
		Email:    in.Email,
		Password: in.Password,
	}
	err = s.dbRepo.InsertUser(ctx, userData)
	if err != nil {
		return nil, err
	}
	token, err = userData.GenerateToken(s.cfg)
	if err != nil {
		return nil, err
	}
	return token, nil
}
