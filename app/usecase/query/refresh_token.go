package query

import (
	"context"
	"github.com/forhomme/app-base/infrastructure/telemetry"
	"github.com/forhomme/app-base/usecase/logger"
	"user-management/app/common/decorator"
	"user-management/app/domain/user"
	"user-management/config"
)

type RefreshTokenHandler decorator.QueryHandler[*user.RefreshToken, *user.Token]

type refreshTokenRepository struct {
	dbRepo user.QueryRepository
	cfg    *config.Config
}

func NewRefreshTokenRepository(dbRepo user.QueryRepository, logger logger.Logger, cfg *config.Config, tracer *telemetry.OtelSdk) decorator.QueryHandler[*user.RefreshToken, *user.Token] {
	return decorator.ApplyQueryDecorators[*user.RefreshToken, *user.Token](
		refreshTokenRepository{dbRepo: dbRepo, cfg: cfg},
		logger,
		tracer,
	)
}

func (r refreshTokenRepository) Handle(ctx context.Context, in *user.RefreshToken) (token *user.Token, err error) {
	userData, err := r.dbRepo.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return nil, err
	}
	err = in.ParsingToken(userData, r.cfg)
	if err != nil {
		return nil, err
	}
	token, err = userData.GenerateToken(r.cfg)
	if err != nil {
		return nil, err
	}
	return token, nil
}
