/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/spf13/cobra"
	"github.com/ts-vis-go/internal/typescript"
)

var MaxDepth int

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
		graph := typescript.Scan(args[0])
		tree := charts.NewTree()
		tree.SetGlobalOptions(charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "100vh",
		}))
		tree.AddSeries("tree", IntoTreeData(graph, MaxDepth), charts.WithTreeOpts(opts.TreeChart{
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

func IntoTreeData(g *typescript.Graph, maxDepth int) []opts.TreeData {
	var treeData []opts.TreeData
	for _, node := range g.NodesByDepth(0) {
		treeData = append(treeData, opts.TreeData{
			Name:     node.Path,
			Children: IntoChildren(g, &node, 1, maxDepth),
		})
	}

	return treeData
}

func IntoChildren(g *typescript.Graph, node *typescript.Node, depth int, maxDepth int) []*opts.TreeData {
	var children []*opts.TreeData
	if depth > maxDepth {
		return children
	}

	edges := g.EdgesFromNode(node)

	for _, edge := range edges {
		children = append(children, &opts.TreeData{
			Name:     edge.To.Path,
			Children: IntoChildren(g, edge.To, depth+1, maxDepth),
		})
	}

	return children
}
