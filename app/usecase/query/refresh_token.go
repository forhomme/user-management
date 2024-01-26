package query

import (
	"context"
	"github.com/forhomme/app-base/usecase/logger"
	"user-management/app/common/decorator"
	"user-management/app/domain/user"
	"user-management/config"
)

type RefreshTokenHandler decorator.QueryHandler[*user.RefreshToken, *user.Token]

type refreshTokenRepository struct {
	dbRepo user.QueryRepository
	logger logger.Logger
	cfg    *config.Config
}

func NewRefreshTokenRepository(dbRepo user.QueryRepository, logger logger.Logger, cfg *config.Config) decorator.QueryHandler[*user.RefreshToken, *user.Token] {
	return decorator.ApplyQueryDecorators[*user.RefreshToken, *user.Token](
		refreshTokenRepository{dbRepo: dbRepo, logger: logger, cfg: cfg},
		logger,
	)
}

func (r refreshTokenRepository) Handle(ctx context.Context, in *user.RefreshToken) (*user.Token, error) {
	userData, err := r.dbRepo.GetUserByEmail(in.Email)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	err = in.ParsingToken(userData, r.cfg)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	token, err := userData.GenerateToken(r.cfg)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	return token, nil
}
