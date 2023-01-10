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
		expect_lineno int32
	}{
		{doc: "", expect_lineno: -1},
		{doc: "i = int(5)", expect_lineno: -1},
		{doc: "envdlib = include(\"https://github.com/tensorchord/envdlib\")\ndef foo():\n    pass\n", expect_lineno: -1},
		{doc: "envdlib = include(\"https://github.com/tensorchord/envdlib\")\ndef build():\n    pass\n", expect_lineno: 1},
		{doc: "def build():\n    pass\n", expect_lineno: 0},
		{doc: "def mod():\n    pass\ndef build():\n    pass\n", expect_lineno: 2},
	}

	for _, tt := range tests {
		t.Run(tt.doc, func(t *testing.T) {
			doc := newDocument(tt.doc)
			node := LoadEnvdEntry(doc)

			if node.IsNull() {
				assert.Equal(t, int32(-1), tt.expect_lineno)
			} else {
				assert.Equal(t, int32(node.StartPoint().Row), tt.expect_lineno)
			}
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
