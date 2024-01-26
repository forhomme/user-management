package query

import (
	"context"
	"errors"
	"github.com/forhomme/app-base/usecase/logger"
	"user-management/app/common/decorator"
	"user-management/app/domain/user"
	"user-management/config"
)

type LoginHandler decorator.QueryHandler[*user.Login, *user.Token]

type loginRepository struct {
	dbRepo user.QueryRepository
	logger logger.Logger
	cfg    *config.Config
}

func NewLoginRepository(dbRepo user.QueryRepository, logger logger.Logger, cfg *config.Config) decorator.QueryHandler[*user.Login, *user.Token] {
	return decorator.ApplyQueryDecorators[*user.Login, *user.Token](
		loginRepository{dbRepo: dbRepo, logger: logger, cfg: cfg},
		logger,
	)
}

func (l loginRepository) Handle(ctx context.Context, in *user.Login) (*user.Token, error) {
	userData, err := l.dbRepo.GetUserByEmail(in.Email)
	if err != nil {
		l.logger.Error(err)
		return nil, err
	}
	if !userData.IsExist() {
		return nil, errors.New("user not found")
	}

	err = in.CheckPassword(userData.Password)
	if err != nil {
		l.logger.Error(err)
		return nil, err
	}

	token, err := userData.GenerateToken(l.cfg)
	if err != nil {
		l.logger.Error(err)
		return nil, err
	}
	return token, nil
}
