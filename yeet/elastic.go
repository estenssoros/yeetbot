package main

import (
	"github.com/spf13/cobra"
)

var elasticCmd = &cobra.Command{
	Use: "elastic",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
