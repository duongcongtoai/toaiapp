package oauth2server

import (
	"toaiapp/registry"

	"github.com/labstack/echo/v4"
)

const ComponentName = "app_oauth2server_component"

var (
	Component *AppOauth2ServerComponent
)

type AppOauth2ServerComponent struct {
	r *registry.Registry
	w int
}

//SetupEcho register routes
func (u *AppOauth2ServerComponent) SetupEcho(e *echo.Echo) error {
	registerRoutes(e)
	return nil
}

func (u *AppOauth2ServerComponent) SetupFromYaml(configFile string) error {
	return nil
}

func (u *AppOauth2ServerComponent) GetWeight() int {
	return u.w
}

func (u *AppOauth2ServerComponent) GetName() string {
	return ComponentName
}
func (u *AppOauth2ServerComponent) Shutdown() error {
	return nil
}

func (u *AppOauth2ServerComponent) SetRegistry(r *registry.Registry) {
	u.r = r
}

func init() {
	Component = &AppOauth2ServerComponent{}
	registry.Instance().Register(ComponentName, Component)
}
