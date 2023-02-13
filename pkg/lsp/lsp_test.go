// Copyright 2023 The envd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lsp

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tensorchord/envd-lsp/pkg/api"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func TestLSP_Hover(t *testing.T) {
	for choice := range api.APIOptions {
		actual := api.APIOptions[choice]
		f := newFixture(t, actual)

		docURI := uri.File("./hover.envd")
		doc := `base(os="ubuntu20.04", language="python3")`

		f.mustWriteDocument("./hover.envd", doc)

		var resp protocol.Hover
		f.mustEditorCall(protocol.MethodTextDocumentHover, protocol.HoverParams{
			TextDocumentPositionParams: protocol.TextDocumentPositionParams{
				TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
				Position: protocol.Position{
					Line:      0,
					Character: 0,
				},
			},
		}, &resp)

		require.Equal(t, resp.Contents.Kind, protocol.Markdown)
		require.NotEmpty(t, resp.Contents.Value)
		require.Contains(t, resp.Contents.Value, "Parameters")
		require.Contains(t, resp.Contents.Value, "\\\n")
	}

}

func TestServer_SemanticToken(t *testing.T) {
	f := newFixture(t, "stable")

	docURI := uri.File("./semantic.envd")
	doc := `
def build():
	pass
def build_GPU():
	pass
`

	f.mustWriteDocument("./semantic.envd", doc)

	var resp protocol.SemanticTokens

	f.mustEditorCall(protocol.MethodSemanticTokensFull, protocol.SemanticTokensParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: docURI},
	}, &resp)

	require.Len(t, resp.Data, 10)
}
