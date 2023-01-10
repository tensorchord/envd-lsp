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

package analysis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tilt-dev/starlark-lsp/pkg/document"
	starlark_query "github.com/tilt-dev/starlark-lsp/pkg/query"
)

func TestLoadEnvdEntry(t *testing.T) {
	tests := []struct {
		doc         string
		expect_data []uint32
	}{
		{doc: "", expect_data: []uint32{}},
		{doc: "i = int(5)", expect_data: []uint32{}},
		{doc: "envdlib = include(\"https://github.com/tensorchord/envdlib\")\ndef foo():\n    pass\n",
			expect_data: []uint32{0, 10, 0, 1, 0}},
		{doc: "envdlib=include(\"https://github.com/tensorchord/envdlib\")\ndef build():\n    pass\n",
			expect_data: []uint32{0, 8, 0, 1, 0, 1, 4, 0, 0, 0}},
		{doc: "def build():\n    pass\n", expect_data: []uint32{0, 4, 0, 0, 0}},
		{doc: "\n\ndef build():\n    pass\n", expect_data: []uint32{2, 4, 0, 0, 0}},
		{doc: "def mod():\n    pass\ndef build():\n    pass\n", expect_data: []uint32{2, 4, 0, 0, 0}},
		{doc: "a = include(\"a.A\")\n\nb = include(\"b.B\")\n", expect_data: []uint32{0, 4, 0, 1, 0, 2, 4, 0, 1, 0}},
		{doc: "a = include(\"a.A\") b = include(\"b.B\")\n", expect_data: []uint32{0, 4, 0, 1, 0, 0, 19, 0, 1, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.doc, func(t *testing.T) {
			ctx := context.Background()
			doc := newDocument(tt.doc)
			analysis, _ := NewAnalyzer(ctx)

			actual_data := analysis.SemanticToken(ctx, doc)
			assert.ElementsMatch(t, actual_data, tt.expect_data)
		})
	}
}

func newDocument(content string) document.Document {
	ctx := context.Background()
	bytes := []byte(content)
	tree, _ := starlark_query.Parse(ctx, bytes)
	doc := document.NewDocument("", bytes, tree)
	return doc
}
