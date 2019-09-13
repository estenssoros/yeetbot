package main

import (
	"fmt"

	"github.com/estenssoros/yeetbot/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	userCmd.AddCommand(userListCmd)
}

var userCmd = &cobra.Command{
	Use: "user",
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "list users in a workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := client.ConfigFromEnv()
		if err != nil {
			return errors.Wrap(err, "client config from env")
		}
		c := config.NewClient(config.Reports[0])
		users, err := c.ListUsers()
		if err != nil {
			return errors.Wrap(err, "client list users")
		}

		for _, user := range users {
			if verbose {
				fmt.Println(user)
			} else {
				fmt.Println(user.ID, user.Name)
			}
		}
		return nil
	},
}
