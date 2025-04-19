package typescript

import (
	"os"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_ts "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

type Node struct {
	Path  string
	Depth int
}

type Edge struct {
	From *Node
	To   *Node
}

type Graph struct {
	Nodes []Node
	Edges []Edge
	MaxDepth int
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

func newGraph() *Graph {
	return &Graph{
		Nodes: []Node{},
		Edges: []Edge{},
		MaxDepth: 20,
	}
}

func Scan(entrypoint string) *Graph {
	parser := tree_sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_ts.LanguageTypescript()))
	resolver := NewResolver(entrypoint)

	graph := newGraph()
	next := Node{Path: entrypoint, Depth: 0}
	graph.Nodes = append(graph.Nodes, next)

	read(parser, resolver, graph, &next)

	return graph
}

func read(parser *tree_sitter.Parser, resolver *Resolver, graph *Graph, node *Node) {
	data, err := os.ReadFile(node.Path)
	if err != nil {
		return
	}

	tree := parser.Parse(data, nil)
	defer tree.Close()

	program := tree.RootNode()
	for _, ref := range references(&data, program) {
		resolved, err := resolver.ResolvePath(ref, node.Path)
		if err != nil {
			continue
		}

		if n := graph.FindNode(resolved); n != nil {
			graph.Edges = append(graph.Edges, Edge{From: node, To: n})
			continue
		}

		next := Node{Path: resolved, Depth: node.Depth+1}
		graph.Nodes = append(graph.Nodes, next)
		graph.Edges = append(graph.Edges, Edge{From: node, To: &next})

		if next.Depth > graph.MaxDepth {
			continue
		}

		read(parser, resolver, graph, &next)
	}
}

func references(src *[]byte, node *tree_sitter.Node) []string {
	ret := []string{}
	cursor := node.Walk()
	defer cursor.Close()

	for _, node := range node.Children(cursor) {
		if node.Kind() != "import_statement" {
			continue
		}

		sourceNode := node.ChildByFieldName("source")
		ret = append(ret, string((*src)[sourceNode.StartByte()+1:sourceNode.EndByte()-1]))
	}

	return ret
}
