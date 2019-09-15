package main

import (
	"fmt"

	"github.com/estenssoros/yeetbot/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var listDmCmd = &cobra.Command{
	Use:   "dm",
	Short: "list direct messages",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := client.ConfigFromEnv()
		if err != nil {
			return errors.Wrap(err, "client config from env")
		}
		c := config.NewClient(config.Reports[0])
		channels, err := c.ListDirectMessageChannels()
		if err != nil {
			return errors.Wrap(err, "client list direct messages")
		}
		for _, channel := range channels {
			fmt.Println(channel.Name, channel.ID, channel.User)
		}
		return nil
	},
}
