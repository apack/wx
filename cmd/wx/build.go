package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/apack/wx/internal/compiler"
	"github.com/apack/wx/internal/generator"
	"github.com/spf13/cobra"
	"golang.org/x/tools/imports"
)

var dir string

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build views",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := compiler.NewCompiler(dir)
		if err != nil {
			return err
		}
		views, err := c.Compile()
		if err != nil {
			return err
		}
		g := generator.NewGenerator(views, dir)
		buf := bytes.NewBuffer(nil)

		err = g.Generate(buf)
		if err != nil {
			return err
		}
		data, err := imports.Process("app.wx.go", buf.Bytes(), nil)
		if err != nil {
			return err
		}
		f, err := os.Create("app.wx.go")
		if err != nil {
			return err
		}
		f.Write(data)
		defer f.Close()
		fmt.Printf("wx: %d views compiled\n", len(views))
		return nil
	},
}

func init() {
	buildCmd.
		PersistentFlags().
		StringVarP(&dir, "dir", "d", ".", "app directory")
}
