package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/ts-vis-go/internal/render"
	"github.com/ts-vis-go/internal/typescript"
	"github.com/ts-vis-go/internal/util/glob"
)

var MaxDepth int
var Filters []string
var Includes []string

var scanCmd = &cobra.Command{
	Use:   "scan",
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
		scanner := typescript.NewScanner(typescript.ScannerOptions{
			Entrypoint: args[0],
		})
		defer scanner.Close()

		graph := scanner.Scan()
		chart := render.BuildGraph(graph, glob.CompileFilters(Includes), glob.CompileFilters(Filters))

		f, _ := os.Create("deps.html")
		chart.Render(f)
	},
}
