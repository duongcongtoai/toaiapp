package main

import (
	"fmt"
	"log"
	"toaiapp/registry"

	"github.com/spf13/cobra"
)

func commandWrapper(callable func(*cobra.Command, []string)) func(*cobra.Command, []string) {
	return func(c *cobra.Command, args []string) {
		fmt.Printf("Debug mode: %t\n", debug)
		fmt.Printf("From config file: %s\n", configFile)
		fmt.Printf("On port: %s\n", listen)
		if err := registry.Instance().SetupFromYaml(configFile, debug); err != nil {
			log.Fatalf("Error :%v", err)
		}
		callable(c, args)
	}
}
