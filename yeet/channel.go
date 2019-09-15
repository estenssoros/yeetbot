package main

import (
	"fmt"

	"github.com/estenssoros/yeetbot/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var listChannelCmd = &cobra.Command{
	Use:   "channel",
	Short: "list channels",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := client.ConfigFromEnv()
		if err != nil {
			return errors.Wrap(err, "client config from env")
		}
		c := config.NewClient(config.Reports[0])
		channels, err := c.ListChannels()
		if err != nil {
			return errors.Wrap(err, "client list channels")
		}
		for _, channel := range channels {
			fmt.Println(channel.ID, channel.Name)
		}
		return nil
	},
}

var channelPurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "purge all messages from a channel",
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
		c := config.NewClient(config.Reports[0])
		messages, err := c.ListMessages(args[0])
		if err != nil {
			return errors.Wrap(err, "client list messages")
		}
		toDelete := map[string]string{}
		for _, message := range messages {
			if err := c.DeleteBotMessage(args[0], message.Ts); err != nil {
				toDelete[message.Ts] = message.Text
			}
		}
		for messageTS := range toDelete {
			if err := c.DeleteUserMessage(args[0], messageTS); err == nil {
				delete(toDelete, messageTS)
			}
		}
		if len(toDelete) > 0 {
			fmt.Println("failed to delete messages")
		}
		for messageTS, messageText := range toDelete {
			fmt.Println(messageTS, messageText)
		}
		return nil
	},
}
