package user

type CommandRepository interface {
	InsertUser(user *User) error
	UpdateUser(id string, updateFn func(u *User) (*User, error)) error
}

type QueryRepository interface {
	GetUserById(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
}
