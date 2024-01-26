package command

import (
	"context"
	"errors"
	"github.com/forhomme/app-base/usecase/logger"
	"user-management/app/common/decorator"
	"user-management/app/domain/user"
)

type ChangePasswordHandler decorator.CommandHandler[*user.ChangePassword]

type changePasswordRepository struct {
	dbRepo user.CommandRepository
	logger logger.Logger
}

func NewChangePasswordRepository(dbRepo user.CommandRepository, logger logger.Logger) decorator.CommandHandler[*user.ChangePassword] {
	return decorator.ApplyCommandDecorators[*user.ChangePassword](
		changePasswordRepository{dbRepo: dbRepo, logger: logger},
		logger,
	)
}

func (c changePasswordRepository) Handle(ctx context.Context, in *user.ChangePassword) error {
	return c.dbRepo.UpdateUser(in.UserId, func(u *user.User) (*user.User, error) {
		if !in.IsValid(u) {
			return nil, errors.New("old or new password is not valid")
		}
		err := u.ChangePassword(in.NewPassword)
		if err != nil {
			return nil, err
		}
		return u, nil
	})
}
