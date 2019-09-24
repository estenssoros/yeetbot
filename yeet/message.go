package main

import (
	"fmt"

	"github.com/estenssoros/yeetbot/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var listMessageCmd = &cobra.Command{
	Use:   "message",
	Short: "list messages in a channel",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("must supply channel argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := client.ConfigFromEnv()
		if err != nil {
			return errors.Wrap(err, "client config from env")
		}
		c := config.NewClient()
		messages, err := c.ListMessages(args[0])
		if err != nil {
			return errors.Wrap(err, "client list messages")
		}
		for _, message := range messages {
			fmt.Println(message)
		}
		return nil
	},
}
