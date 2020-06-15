package main

import (
	"log"
	"toaiapp/registry"
	"os"

	_ "toaiapp/app/auth"
	_ "toaiapp/auth"
	_ "toaiapp/auth/db/postgresql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

var (
	debug bool

	listen      string
	configFile  string
	rootCommand = &cobra.Command{

		Use:   "toai",
		Short: "Test Api root Command",
		Long:  ``,
	}
	serveCommand = &cobra.Command{
		Use:   "serve",
		Short: "Run the web server",
		Long:  ``,
		Run:   commandWrapper(serve),
	}
)

func serve(cmd *cobra.Command, args []string) {
	e := echo.New()
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

func init() {
	rootCommand.AddCommand(serveCommand)
	rootCommand.AddCommand(authCommand)
	rootCommand.PersistentFlags().StringVar(&listen, "listen", "", "")
	rootCommand.PersistentFlags().StringVar(&configFile, "configFile", "", "")
	rootCommand.PersistentFlags().BoolVar(&debug, "debug", true, "")
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(-1)
	}
}
