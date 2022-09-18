package main

import (
	"github.com/apack/wx"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize application",
	RunE: func(cmd *cobra.Command, args []string) error {
		return wx.Initialize("app")
	},
}
