package auth

import (
	"github.com/labstack/echo/v4"
)

type AuthDB interface {
	SetupEcho(e *echo.Echo) error
	SetConfig(Configuration)

	FindUserByID(id float64) (*User, error)
	FindUserByName(name string) (*User, error)
	CreateUser(*User) error

	Initialize() error
	DB() (AuthDB, error)
	FromContext(echo.Context) (AuthDB, error)
}
