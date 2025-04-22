package render

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
	"github.com/ts-vis-go/internal/model"
	"github.com/ts-vis-go/internal/util/glob"
)

func BuildGraph(g *model.Graph, includes []glob.Glob, filters []glob.Glob) render.Renderer {
	chart := charts.NewGraph()
	chart.SetGlobalOptions(charts.WithInitializationOpts(opts.Initialization{
		Width:  "100%",
		Height: "100vh",
	}))

	nodes, links := mapNodesAndLinks(g, includes, filters)

	chart.AddSeries("graph", nodes, links, charts.WithGraphChartOpts(opts.GraphChart{
		Layout: "force",
		Force: &opts.GraphForce{
			Repulsion: 100,
			Gravity: 0.1,
		},
		Draggable: opts.Bool(true),
		Roam: opts.Bool(true),
	}))

	return chart
}

func mapNodesAndLinks(g *model.Graph, includes []glob.Glob, filters []glob.Glob) ([]opts.GraphNode, []opts.GraphLink) {
	var nodes []opts.GraphNode
	var links []opts.GraphLink

	for _, node := range g.Nodes {
		if !glob.AnyMatches(includes, node.Path) {
			continue
		}

		if glob.AnyMatches(filters, node.Path) {
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
		if !glob.AnyMatches(includes, edge.From.Path, edge.To.Path) {
			continue
		}

		if glob.AnyMatches(filters, edge.From.Path, edge.To.Path) {
			continue
		}

		links = append(links, opts.GraphLink{
			Source: edge.From.Name,
			Target: edge.To.Name,
		})
	}

	return nodes, links
}


func nodesContainNode(nodes []opts.GraphNode, name any) bool {
	for _, node := range nodes {
		if node.Name == name {
			return true
		}
	}

	return false
}
