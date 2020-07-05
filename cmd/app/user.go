package main

import (
	"fmt"
	"os"

	"toaiapp/auth"

	"github.com/spf13/cobra"
)

var authCommand = &cobra.Command{
	Use:   "auth",
	Short: "Authentication",
	Long:  ``,
}

var addUserCommand = &cobra.Command{
	Use:   "add <username> <password>",
	Short: "Add a user",
	Long:  ``,
	Run: commandWrapper(func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Error: Give username and password")
			os.Exit(1)
		}
		user := &auth.User{
			Name: args[0],
		}
		user.SetPassword(args[1])
		db := auth.GetDB()

		if err := db.CreateUser(user); err != nil {

			fmt.Printf("Error: Error creating user: %v", err)
			os.Exit(1)
		}
		fmt.Printf("User '%s' successfully created.\n", user.Name)
	}),
}

func init() {
	authCommand.AddCommand(addUserCommand)
}
