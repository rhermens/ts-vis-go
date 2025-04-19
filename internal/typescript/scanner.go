package typescript

import (
	"os"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_ts "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
	"github.com/ts-vis-go/internal/model"
)

func Scan(entrypoint string) *model.Graph {
	parser := tree_sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_ts.LanguageTypescript()))
	resolver := NewResolver(entrypoint)

	graph := model.NewGraph()
	next := model.Node{Path: entrypoint, Depth: 0}
	graph.Nodes = append(graph.Nodes, next)

	read(parser, resolver, graph, &next)

	return graph
}

func read(parser *tree_sitter.Parser, resolver *Resolver, graph *model.Graph, node *model.Node) {
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
			graph.Edges = append(graph.Edges, model.Edge{From: node, To: n})
			continue
		}

		next := model.Node{Path: resolved, Depth: node.Depth+1}
		graph.Nodes = append(graph.Nodes, next)
		graph.Edges = append(graph.Edges, model.Edge{From: node, To: &next})

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
