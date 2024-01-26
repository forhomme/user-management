package query

import (
	"context"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/google/uuid"
	"user-management/app/common/decorator"
	"user-management/app/domain/user"
	"user-management/config"
)

type SignUpHandler decorator.QueryHandler[*user.SignUp, *user.Token]

type signUpRepository struct {
	dbRepo user.CommandRepository
	logger logger.Logger
	cfg    *config.Config
}

func NewSignUpRepository(dbRepo user.CommandRepository, logger logger.Logger, cfg *config.Config) decorator.QueryHandler[*user.SignUp, *user.Token] {
	return decorator.ApplyQueryDecorators[*user.SignUp, *user.Token](
		signUpRepository{dbRepo: dbRepo, logger: logger, cfg: cfg},
		logger,
	)
}

func (s signUpRepository) Handle(ctx context.Context, in *user.SignUp) (*user.Token, error) {
	err := in.HashPassword()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	userData := &user.User{
		RoleId:   in.RoleId,
		UserId:   uuid.New().String(),
		Email:    in.Email,
		Password: in.Password,
	}
	err = s.dbRepo.InsertUser(userData)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	token, err := userData.GenerateToken(s.cfg)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return token, nil
}
