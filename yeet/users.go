package main

import (
	"fmt"

	"github.com/estenssoros/yeetbot/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var listUserCmd = &cobra.Command{
	Use:   "user",
	Short: "list users in a workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := client.ConfigFromEnv()
		if err != nil {
			return errors.Wrap(err, "client config from env")
		}
		c := config.NewClient()
		users, err := c.ListUsers()
		if err != nil {
			return errors.Wrap(err, "client list users")
		}

		for _, user := range users {
			fmt.Println(user.ID, user.Name)
		}
		return nil
	},
}
