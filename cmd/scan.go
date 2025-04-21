package cmd

import (
	"errors"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
	"github.com/ts-vis-go/internal/model"
	"github.com/ts-vis-go/internal/typescript"
)

var MaxDepth int
var Filters []string

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
		chart := chartGraph(graph, compileFilters(Filters))

		f, _ := os.Create("deps.html")
		chart.Render(f)
	},
}

func mapNodesAndLinks(g *model.Graph, filters []glob.Glob) ([]opts.GraphNode, []opts.GraphLink) {
	var nodes []opts.GraphNode
	var links []opts.GraphLink

	for _, node := range g.Nodes {
		if !anyMatches(filters, node.Path) {
			continue
		}

		if nodesContainNode(nodes, node.Name) {
			continue
		}

		nodes = append(nodes, opts.GraphNode{
			Name: node.Name,
		})
	}

	for _, edge := range g.Edges {
		if !anyMatches(filters, edge.From.Path) || !anyMatches(filters, edge.To.Path) {
			continue
		}

		links = append(links, opts.GraphLink{
			Source: edge.From.Name,
			Target: edge.To.Name,
		})
	}

	return nodes, links
}

func chartGraph(g *model.Graph, filters []glob.Glob) render.Renderer {
	chart := charts.NewGraph()
	chart.SetGlobalOptions(charts.WithInitializationOpts(opts.Initialization{
		Width:  "100%",
		Height: "100vh",
	}))

	nodes, links := mapNodesAndLinks(g, filters)

	chart.AddSeries("graph", nodes, links, charts.WithGraphChartOpts(opts.GraphChart{
		Layout: "force",
		Force: &opts.GraphForce{
			Repulsion: 100,
		},
		Roam: opts.Bool(true),
	}))

	return chart
}

func anyMatches(filters []glob.Glob, str string) bool {
	for _, filter := range filters {
		if filter.Match(str) {
			return true
		}
	}

	return false
}

func compileFilters(input []string) []glob.Glob {
	var filters []glob.Glob
	for _, filter := range input {
		filters = append(filters, glob.MustCompile(filter))
	}

	return filters
}

func nodesContainNode(nodes []opts.GraphNode, name any) bool {
	for _, node := range nodes {
		if node.Name == name {
			return true
		}
	}

	return false
}
