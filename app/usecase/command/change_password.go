package command

import (
	"context"
	"errors"
	"github.com/forhomme/app-base/infrastructure/telemetry"
	"github.com/forhomme/app-base/usecase/logger"
	"user-management/app/common/decorator"
	"user-management/app/domain/user"
)

type ChangePasswordHandler decorator.CommandHandler[*user.ChangePassword]

type changePasswordRepository struct {
	dbRepo user.CommandRepository
}

func NewChangePasswordRepository(dbRepo user.CommandRepository, logger logger.Logger, tracer *telemetry.OtelSdk) decorator.CommandHandler[*user.ChangePassword] {
	return decorator.ApplyCommandDecorators[*user.ChangePassword](
		changePasswordRepository{dbRepo: dbRepo},
		logger,
		tracer,
	)
}

func (c changePasswordRepository) Handle(ctx context.Context, in *user.ChangePassword) error {
	return c.dbRepo.UpdateUser(ctx, in.UserId, func(u *user.User) (out *user.User, err error) {
		if !in.IsValid(u) {
			err = errors.New("old or new password is not valid")
			return nil, err
		}
		err = u.ChangePassword(in.NewPassword)
		if err != nil {
			return nil, err
		}
		return u, nil
	})
}
