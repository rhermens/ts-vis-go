package model

import "path/filepath"

type Node struct {
	Name  string
	Path  string
	Depth int
}

type Edge struct {
	From *Node
	To   *Node
}

type Graph struct {
	Nodes    []Node
	Edges    []Edge
	MaxDepth int
}

func NewNode(path string, depth int) Node {
	return Node{Path: path, Depth: 0, Name: filepath.Base(path)}
}

func (g *Graph) FindNode(path string) *Node {
	for _, node := range g.Nodes {
		if node.Path == path {
			return &node
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

func (g *Graph) NodesByDepth(depth int) []Node {
	nodes := []Node{}
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

func NewGraph() *Graph {
	return &Graph{
		Nodes:    []Node{},
		Edges:    []Edge{},
		MaxDepth: 20,
	}
}

func (g *Graph) AddNode(node Node) {
	g.Nodes = append(g.Nodes, node)
}

func (g *Graph) AddEdge(from *Node, to *Node) {
	g.Edges = append(g.Edges, Edge{From: from, To: to})
}
