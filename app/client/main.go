package main

import (
	"context"
	"net/http"
	_ "toaiapp/app/auth"
	_ "toaiapp/auth"
	_ "toaiapp/auth/db/postgresql"

	"github.com/labstack/echo"
	"golang.org/x/oauth2"
)

var (
	config = oauth2.Config{
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

func main() {
	e := echo.New()
	e.GET("/", home)
	e.GET("/oauth2", authorize)
	e.Start(":8084")
}

func home(c echo.Context) error {
	u := config.AuthCodeURL("xyz")
	return c.Redirect(http.StatusFound, u)
}

func authorize(c echo.Context) error {
	state := c.FormValue("state")
	if state != "xyz" {
		c.JSON(http.StatusNotFound, map[string]string{"message": "State invalid"})
	}

	code := c.FormValue("code")
	if code == "" {
		c.JSON(http.StatusNotFound, map[string]string{"message": "Code not found"})
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusFound, token)
}
