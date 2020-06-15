package auth

import (
	"toaiapp/registry"

	"github.com/labstack/echo/v4"
)

const ComponentName = "app_auth_component"

var (
	Component *AppAuthComponent
)

type AppAuthComponent struct {
	r *registry.Registry
	w int
}

//SetupEcho register routes
func (u *AppAuthComponent) SetupEcho(e *echo.Echo) error {
	apiRegisterRoutes(e)
	return nil
}

func (u *AppAuthComponent) SetupFromYaml(configFile string, debug bool) error {
	return nil
}

func (u *AppAuthComponent) GetWeight() int {
	return u.w
}

func (u *AppAuthComponent) GetName() string {
	return "App_Auth"
}
func (u *AppAuthComponent) Shutdown() error {
	return nil
}

func (u *AppAuthComponent) SetRegistry(r *registry.Registry) {
	u.r = r
}

func init() {
	Component = &AppAuthComponent{}
	registry.Instance().Register(ComponentName, Component)
}
