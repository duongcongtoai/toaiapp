package auth

import (
	"context"
	"fmt"
	"net/http"
	"toaiapp/auth"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

var (
	oauth2Config = oauth2.Config{
		ClientID:     "client_app_id",
		ClientSecret: "client_secret",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:8084/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:   "http://localhost:8082/oauth/authorize",
			TokenURL:  "http://app:8082/oauth/get_token",
			AuthStyle: oauth2.AuthStyleAutoDetect,
		},
	}
)

func registerRoutes(e *echo.Echo) {
	e.GET("/", home)
	e.GET("/oauth2", authorize)

	// gr := e.Group("/api/v1/user")

	// gr.POST("/login/toaiapp", auth.AuthorizationWrapper(loginWithToaiApp, ""))
	// gr.Use(auth.MiddlewareSessionAuth())
	// gr.GET("/", auth.AuthorizationWrapper(getUser, auth.UserGet))
}

func home(c echo.Context) error {
	u := oauth2Config.AuthCodeURL("xyz")
	return c.Redirect(http.StatusFound, u)
}

func authorize(c echo.Context) error {
	state := c.FormValue("state")
	if state != "xyz" {

		return c.JSON(http.StatusNotFound, map[string]string{"message": "State invalid"})
	}
	code := c.FormValue("code")
	if code == "" {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Code not found"})
	}
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusFound, token)
}

type ResultGet struct {
	Greeting string
}

func loginWithToaiApp(c echo.Context, u *auth.User) error {
	return nil
}

func getUser(c echo.Context, u *auth.User) error {
	if u != auth.Component.UserGuest {
		return c.JSON(http.StatusOK, ResultGet{fmt.Sprintf("Hello there, %s", u.Name)})
	}
	return c.JSON(http.StatusOK, ResultGet{"Hello there guest"})
}

func sendToken(c echo.Context, u *auth.User) error {
	token, err := u.GenerateToken()
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
