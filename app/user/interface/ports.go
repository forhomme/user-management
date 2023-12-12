package _interface

import (
	"github.com/forhomme/app-base/usecase/controller"
)

type ControllerPorts interface {
	SignUp(c controller.Context) (err error)
	Login(c controller.Context) (err error)
	RefreshToken(c controller.Context) (err error)
	ChangePassword(c controller.Context) (err error)
	GetMenu(c controller.Context) (err error)
}
