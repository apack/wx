package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version of wx",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wx version: v0.1")
	},
}
