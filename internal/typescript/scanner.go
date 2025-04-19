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
	nodes []Node
	edges []Edge
}

func (g *Graph) ContainsNode(path string) bool {
	for _, node := range g.nodes {
		if node.Path == path {
			return true
		}
	}

	return false
}

func Scan(entrypoint string) *Graph {
	parser := tree_sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_ts.LanguageTypescript()))
	resolver := NewResolver(entrypoint)

	graph := &Graph{}

	read(parser, resolver, graph, &Node{Path: entrypoint, Depth: 0})

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

		next := Node{Path: resolved, Depth: node.Depth + 1}
		if graph.ContainsNode(next.Path) {
			return
		}

		graph.nodes = append(graph.nodes, next)
		graph.edges = append(graph.edges, Edge{From: node, To: &next})
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
