package model

import (
	"log/slog"
	"strings"
)

type Node struct {
	Index int
	Name  string
	Path  string
	Depth int
}

type Edge struct {
	From *Node
	To   *Node
}

type Graph struct {
	Cwd      string
	Nodes    []*Node
	Edges    []Edge
	MaxDepth int
}

func NewNode(path string, depth int) *Node {
	return &Node{Path: path, Depth: 0}
}

func (g *Graph) FindNode(path string) *Node {
	for _, node := range g.Nodes {
		if node.Path == path {
			return node
		}
	}

	return nil
}

func (g *Graph) ContainsNode(path string) bool {
	if g.FindNode(path) != nil {
		return true
	}

	return false
}

func (g *Graph) NodesByDepth(depth int) []*Node {
	nodes := []*Node{}
	for _, node := range g.Nodes {
		if node.Depth == depth {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (g *Graph) EdgesFromNode(node *Node) []Edge {
	edges := []Edge{}
	for _, edge := range g.Edges {
		if edge.From.Path == node.Path {
			edges = append(edges, edge)
		}
	}

	return edges
}

func NewGraph(cwd string) *Graph {
	return &Graph{
		Cwd:      cwd,
		Nodes:    []*Node{},
		Edges:    []Edge{},
		MaxDepth: 20,
	}
}

func (g *Graph) AddNode(node *Node) {
	node.Index = len(g.Nodes)
	node.Name = strings.TrimPrefix(node.Path, g.Cwd)
	g.Nodes = append(g.Nodes, node)
	slog.Info("Added node ", "n", len(g.Nodes)+1, "name", node.Name)
}

func (g *Graph) AddEdge(from *Node, to *Node) {
	g.Edges = append(g.Edges, Edge{From: from, To: to})
	slog.Info("Added edge ", "n", len(g.Edges)+1, "from", from.Name, "to", to.Name)
}
