package typescript

import (
	"log/slog"
	"os"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_ts "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
	"github.com/ts-vis-go/internal/model"
)

type Scanner struct {
	Parser   *tree_sitter.Parser
	Resolver *Resolver
	Options  ScannerOptions
}

type ScannerOptions struct {
	Entrypoint string
}

func NewScanner(options ScannerOptions) *Scanner {
	parser := tree_sitter.NewParser()

	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_ts.LanguageTypescript()))
	resolver := NewResolver(options.Entrypoint)

	return &Scanner{
		Parser:   parser,
		Resolver: resolver,
		Options:  options,
	}
}

func (s *Scanner) Close() {
	s.Parser.Close()
}

func (s *Scanner) Scan() *model.Graph {
	graph := model.NewGraph(s.Resolver.cwd)

	next := model.NewNode(s.Options.Entrypoint, 0)
	graph.AddNode(next)
	s.next(graph, next)

	return graph
}

func (s *Scanner) next(graph *model.Graph, node *model.Node) {
	slog.Debug("next", "path", node.Path)

	data, err := os.ReadFile(node.Path)
	if err != nil {
		slog.Error("Error reading file", "path", node.Path, "error", err)
		return
	}

	tree := s.Parser.Parse(data, nil)
	defer tree.Close()

	program := tree.RootNode()
	for _, ref := range references(&data, program) {
		resolved, err := s.Resolver.ResolvePath(ref, node.Path)
		if err != nil {
			slog.Debug("Error resolving path", "path", ref, "error", err)
			continue
		}
		slog.Debug("Resolved", "reference", ref, "resolved", resolved)

		if n := graph.FindNode(resolved); n != nil {
			graph.AddEdge(node, n)
			continue
		}

		next := model.NewNode(resolved, node.Depth+1)
		graph.AddNode(next)
		graph.AddEdge(node, next)

		s.next(graph, next)
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
