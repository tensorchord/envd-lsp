// Copyright 2023 The envd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package query

import (
	"strings"

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
func LoadEnvdEntries(doc starlark_query.DocumentContent) []*sitter.Node {
	nodes := []*sitter.Node{}
	starlark_query.Query(doc.Tree().RootNode(), `(function_definition) @function-name`,
		func(q *sitter.Query, match *sitter.QueryMatch) bool {
			for _, c := range match.Captures {
				id := c.Node.ChildByFieldName("name")
				curFuncName := doc.Content(id)
				if strings.Contains(curFuncName, string(SemanticTokenBuild)) {
					nodes = append(nodes, id)
					break
				}
			}
			return true
		})
	return nodes
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
