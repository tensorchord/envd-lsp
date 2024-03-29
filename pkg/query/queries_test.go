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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tilt-dev/starlark-lsp/pkg/document"
	starlark_query "github.com/tilt-dev/starlark-lsp/pkg/query"
)

func TestLoadEnvdEntry(t *testing.T) {
	tests := []struct {
		doc           string
		expect_lineno []uint32
		uint32
	}{
		{doc: "", expect_lineno: []uint32{}},
		{doc: "i = int(5)", expect_lineno: []uint32{}},
		{doc: "envdlib = include(\"https://github.com/tensorchord/envdlib\")\ndef foo():\n    pass\n", expect_lineno: []uint32{}},
		{doc: "envdlib = include(\"https://github.com/tensorchord/envdlib\")\ndef build():\n    pass\n", expect_lineno: []uint32{1}},
		{doc: "def build():\n    pass\n", expect_lineno: []uint32{0}},
		{doc: "def mod():\n    pass\ndef build():\n    pass\n", expect_lineno: []uint32{2}},
		{doc: "def build_gpu():\n    pass\ndef build():\n    pass\n", expect_lineno: []uint32{0, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.doc, func(t *testing.T) {
			doc := newDocument(tt.doc)
			nodes := LoadEnvdEntries(doc)

			actual := []uint32{}
			for _, node := range nodes {
				actual = append(actual, node.StartPoint().Row)
			}
			assert.ElementsMatch(t, actual, tt.expect_lineno)
		})
	}
}

func TestLoadModuleImport(t *testing.T) {
	tests := []struct {
		doc           string
		expect_lineno []uint32
	}{
		{doc: "", expect_lineno: []uint32{}},
		{doc: "i = int(5)", expect_lineno: []uint32{}},
		{doc: "envdlib = include(\"https://github.com/tensorchord/envdlib\")", expect_lineno: []uint32{0}},
		{doc: "envdlib=include(\"https://github.com/tensorchord/envdlib\")", expect_lineno: []uint32{0}},
		{doc: "a = include(\"a.A\")\n\nb = include(\"b.B\")\n", expect_lineno: []uint32{0, 2}},
		{doc: "a = include(\"a.A\") b = include(\"b.B\")\n", expect_lineno: []uint32{0, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.doc, func(t *testing.T) {
			doc := newDocument(tt.doc)
			nodes := LoadModuleImport(doc)

			actual := []uint32{}
			for _, node := range nodes {
				actual = append(actual, node.StartPoint().Row)
			}
			assert.ElementsMatch(t, actual, tt.expect_lineno)
		})
	}
}

func newDocument(content string) starlark_query.DocumentContent {
	ctx := context.Background()
	bytes := []byte(content)
	tree, _ := starlark_query.Parse(ctx, bytes)
	doc := document.NewDocument("", bytes, tree)
	return doc
}
