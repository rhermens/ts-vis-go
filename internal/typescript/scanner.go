package typescript

import (
	"fmt"
	"os"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_ts "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

func Scan(entrypoint string) {
	parser := tree_sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_ts.LanguageTypescript()))
	resolver := NewResolver(entrypoint)

	Read(parser, resolver, entrypoint)
}

func Read(parser *tree_sitter.Parser, resolver *Resolver, entrypoint string) {
	fmt.Printf("Reading %s\n", entrypoint)

	data, err := os.ReadFile(entrypoint)
	if err != nil {
		return
	}

	tree := parser.Parse(data, nil)
	defer tree.Close()

	program := tree.RootNode()
	for _, ref := range References(&data, program) {
		resolved, err := resolver.ResolvePath(ref, entrypoint)
		if err != nil {
			continue
		}

		Read(parser, resolver, resolved);
	}
}

func References(src *[]byte, node *tree_sitter.Node) []string {
	ret := []string{}
	cursor := node.Walk()
	defer cursor.Close()

	for _, node := range node.Children(cursor) {
		if node.Kind() != "import_statement" {
			continue
		}

		sourceNode := node.ChildByFieldName("source")
		ret = append(ret, string((*src)[sourceNode.StartByte() + 1:sourceNode.EndByte() - 1]))
	}

	return ret
} 
