package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
}

func init() {
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(channelCmd)
	rootCmd.AddCommand(messageCmd)
	rootCmd.AddCommand(dmCmd)
	rootCmd.AddCommand(elasticCmd)
}

var rootCmd = &cobra.Command{
	Use:   "yeet",
	Short: "Yeet yeet!",
}

func main() {
	start := time.Now()
	defer func() {
		logrus.Infof("process took %v", time.Since(start))
	}()
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
