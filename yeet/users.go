package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	usersCmd.AddCommand(usersListCmd)
}

var usersCmd = &cobra.Command{
	Use: "users",
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "list users in a workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newClient()
		users, err := client.ListUsers()
		if err != nil {
			return errors.Wrap(err, "client list users")
		}
		for _, u := range users {
			fmt.Println(u.ID, u.Name)
		}
		return nil
	},
}
