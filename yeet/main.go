package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show more text")
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(channelCmd)
	rootCmd.AddCommand(messageCmd)
	rootCmd.AddCommand(dmCmd)
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
