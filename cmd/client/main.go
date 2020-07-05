package main

import (
	"fmt"
	"log"
	"os"
	_ "toaiapp/auth"
	_ "toaiapp/auth/db/postgresql"
	_ "toaiapp/auth/session/psql"
	_ "toaiapp/cmd/client/auth"
	"toaiapp/registry"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
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
	serveCommand = &cobra.Command{
		Use:   "serve",
		Short: "serve oauth2 client app server",
		Long:  ``,
		Run:   parseConfig(serve),
	}
	listen     string
	configFile string
)

func serve(cmd *cobra.Command, args []string) {
	e := echo.New()

	// e.GET("/", home)
	// e.GET("/oauth2", authorize)
	// e.Start(":8084")
	// e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	if err := registry.Instance().SetupEcho(e); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	err := e.Start(listen)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	if err := registry.Instance().Shutdown(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}
func parseConfig(callable func(*cobra.Command, []string)) func(c *cobra.Command, args []string) {
	return func(c *cobra.Command, args []string) {
		fmt.Printf("From config file: %s\n", configFile)
		fmt.Printf("On port: %s\n", listen)
		if err := registry.Instance().SetupFromYaml(configFile); err != nil {
			log.Fatalf("Error :%v", err)
		}
		callable(c, args)
	}
}

func main() {
	// serveCommand.AddCommand(serveCommand)
	// serveCommand.AddCommand(authCommand)
	serveCommand.PersistentFlags().StringVar(&listen, "listen", "", "")
	serveCommand.PersistentFlags().StringVar(&configFile, "configFile", "", "")
	if err := serveCommand.Execute(); err != nil {
		os.Exit(-1)
	}
}

// func home(c echo.Context) error {
// 	u := config.AuthCodeURL("xyz")
// 	return c.Redirect(http.StatusFound, u)
// }

// func authorize(c echo.Context) error {
// 	state := c.FormValue("state")
// 	if state != "xyz" {
// 		return c.JSON(http.StatusNotFound, map[string]string{"message": "State invalid"})
// 	}

// 	code := c.FormValue("code")
// 	if code == "" {
// 		return c.JSON(http.StatusNotFound, map[string]string{"message": "Code not found"})
// 	}

// 	token, err := config.Exchange(context.Background(), code)
// 	if err != nil {
// 		return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
// 	}
// 	return c.JSON(http.StatusFound, token)
// }

// func redirectAuthen(c echo.Context) error {
// }
