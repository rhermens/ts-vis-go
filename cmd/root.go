package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ts-vis",
	Short: "Graph viewer of typescript dependencies",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().IntVarP(&MaxDepth, "max-depth", "d", 500, "Maximum depth of the graph to graph.")
	scanCmd.Flags().StringArrayVarP(&Filters, "filter", "f", []string{"**"}, "Filter file names")
}

