package user

import "context"

type CommandRepository interface {
	InsertUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, id string, updateFn func(u *User) (*User, error)) error
}

type QueryRepository interface {
	GetUserById(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
