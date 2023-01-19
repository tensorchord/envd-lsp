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

		require.NotEmpty(t, resp.Contents.Value)
	}

}

func TestLSP_Hover_Bad_API(t *testing.T) {
	f := newFixture(t, "bad-api")

	docURI := uri.File("./bad-api.envd")
	doc := `base(os="ubuntu20.04", language="python3")`

	f.mustWriteDocument("./bad-api.envd", doc)

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

	require.Empty(t, resp.Contents.Value)
}
