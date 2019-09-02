package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	slackurl = "https://slack.com/api"
	token    = "xoxb-708948424145-745505481845-LcRoch5r9iL94t8gpH4m7dhz"
)

func init() {
	rootCmd.AddCommand(usersCmd)
}

var rootCmd = &cobra.Command{
	Use:   "yeet",
	Short: "Yeet yeet!",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
