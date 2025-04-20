package cmd

import (
	"errors"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
	"github.com/ts-vis-go/internal/model"
	"github.com/ts-vis-go/internal/typescript"
)

var MaxDepth int
var Filter []string

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
			Filters: compileFilters(),
		})
		defer scanner.Close()

		graph := scanner.Scan()
		tree := charts.NewTree()
		tree.SetGlobalOptions(charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "100vh",
		}))
		tree.AddSeries(args[0], intoTreeData(graph, MaxDepth), charts.WithTreeOpts(opts.TreeChart{
			Orient: "LR",
			Roam:   opts.Bool(true),
			Left:   "0",
			Right:  "0",
			Top:    "0",
			Bottom: "0",
		}))

		f, _ := os.Create("deps.html")
		tree.Render(f)
	},
}

func compileFilters() []glob.Glob {
	var filters []glob.Glob
	for _, filter := range Filter {
		filters = append(filters, glob.MustCompile(filter))
	}

	return filters
}

func intoTreeData(g *model.Graph, maxDepth int) []opts.TreeData {
	var treeData []opts.TreeData
	for _, node := range g.NodesByDepth(0) {
		treeData = append(treeData, opts.TreeData{
			Name:     node.Name,
			Children: intoChildren(g, &node, 1, maxDepth),
		})
	}

	return treeData
}

func intoChildren(g *model.Graph, node *model.Node, depth int, maxDepth int) []*opts.TreeData {
	var children []*opts.TreeData
	if depth > maxDepth {
		return children
	}

	edges := g.EdgesFromNode(node)

	for _, edge := range edges {
		children = append(children, &opts.TreeData{
			Name:     edge.To.Name,
			Children: intoChildren(g, edge.To, depth+1, maxDepth),
		})
	}

	return children
}
