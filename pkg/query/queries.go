package query

import (
	sitter "github.com/smacker/go-tree-sitter"
	starlark_query "github.com/tilt-dev/starlark-lsp/pkg/query"
	"go.lsp.dev/protocol"
)

const (
	SemanticTokenBuild  protocol.SemanticTokenTypes = "build"
	SemanticTokenImport protocol.SemanticTokenTypes = "include"
)

var TokenTypes = []protocol.SemanticTokenTypes{SemanticTokenBuild, SemanticTokenImport}

// Catch the semantic token of `build` inside `def build()`, as envd build file entry
//
// As it will located at a line, only lineno is required actually
func LoadEnvdEntry(doc starlark_query.DocumentContent) *sitter.Node {
	node := sitter.Node{}
	starlark_query.Query(doc.Tree().RootNode(), `(function_definition) @function-name`,
		func(q *sitter.Query, match *sitter.QueryMatch) bool {
			for _, c := range match.Captures {
				id := c.Node.ChildByFieldName("name")
				curFuncName := doc.Content(id)
				if curFuncName == string(SemanticTokenBuild) {
					node = *id
					break
				}
			}
			return true
		})
	return &node
}

// Catch the semantic token of `include` such as `envdlib = include("https://github.com/tensorchord/envdlib")`,
// as envd build file module import
//
// As it will located at a line, only lineno is required actually
func LoadModuleImport(doc starlark_query.DocumentContent) []*sitter.Node {
	nodes := []*sitter.Node{}
	starlark_query.Query(doc.Tree().RootNode(), `(call) @call`,
		func(q *sitter.Query, match *sitter.QueryMatch) bool {
			for _, c := range match.Captures {
				id := c.Node.ChildByFieldName("function")
				curFuncName := doc.Content(id)
				if curFuncName == string(SemanticTokenImport) {
					nodes = append(nodes, id)
					break
				}
			}
			return true
		})
	return nodes
}
