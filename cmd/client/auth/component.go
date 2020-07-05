package auth

import (
	"toaiapp/registry"

	"github.com/labstack/echo/v4"
)

const ComponentName = "client_auth_component"

var (
	Component *ClientAuthComponent
)

type ClientAuthComponent struct {
	r *registry.Registry
	w int
}

//SetupEcho register routes
func (u *ClientAuthComponent) SetupEcho(e *echo.Echo) error {
	registerRoutes(e)
	return nil
}

func (u *ClientAuthComponent) SetupFromYaml(configFile string) error {
	return nil
}

func (u *ClientAuthComponent) GetWeight() int {
	return u.w
}

func (u *ClientAuthComponent) GetName() string {
	return ComponentName
}
func (u *ClientAuthComponent) Shutdown() error {
	return nil
}

func (u *ClientAuthComponent) SetRegistry(r *registry.Registry) {
	u.r = r
}

func init() {
	Component = &ClientAuthComponent{}
	registry.Instance().Register(ComponentName, Component)
}
