package main

import "github.com/spf13/cobra"

func init() {
	slackCmd.AddCommand(listCmd)
}

var slackCmd = &cobra.Command{
	Use:   "slack",
	Short: "do slack stuff",
}

func init() {
	listCmd.AddCommand(listUserCmd)
	listCmd.AddCommand(listDmCmd)
	listCmd.AddCommand(listChannelCmd)
	listCmd.AddCommand(listMessageCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list operations",
}
