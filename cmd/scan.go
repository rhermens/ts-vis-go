/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"github.com/ts-vis-go/internal/typescript"
)

var scanCmd = &cobra.Command{
	Use: "scan",
	Short: "Scan entrypoint for references",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("exactly one argument is required")
		}

		if _, err := os.Stat(args[0]); err != nil {
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		graph := typescript.Scan(args[0])
		spew.Dump(graph)
	},
}
