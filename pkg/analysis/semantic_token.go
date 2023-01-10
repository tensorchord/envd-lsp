package analysis

import (
	"context"
	"sort"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"

	"github.com/tilt-dev/starlark-lsp/pkg/document"

	"github.com/tensorchord/envd-lsp/pkg/query"
)

func (a *Analyzer) SemanticToken(ctx context.Context, doc document.Document) []uint32 {
	logger := protocol.LoggerFromContext(ctx)
	logger.Debug("serve semantic token request", zap.String("path", string(doc.URI())))

	entry := query.LoadEnvdEntry(doc)
	imports := query.LoadModuleImport(doc)

	// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocument_semanticTokens
	// absolute semantic stack view
	// [{line:2, startChar:5, length:3, tokenType:0, tokenModifiers:3}] -> [[2,5,3,0,3]]
	// tokenType 0 -> SemanticTokenBuild
	// tokenType 1 -> SemanticTokenImport
	absolute_view := [][5]uint32{}

	if !entry.IsNull() && entry.StartPoint().Row == entry.EndPoint().Row {
		absolute_view = append(absolute_view, [5]uint32{entry.StartPoint().Row, entry.StartPoint().Column,
			entry.EndPoint().Row - entry.StartPoint().Row, 0, 0})
	} else if !entry.IsNull() && entry.StartPoint().Row != entry.EndPoint().Row {
		logger.Warn(`envd entry "def build()" not in single line`, zap.Uint32("start", entry.StartPoint().Row),
			zap.Uint32("end", entry.EndPoint().Row))
	}

	for _, item := range imports {
		if !item.IsNull() && item.StartPoint().Row == item.EndPoint().Row {
			absolute_view = append(absolute_view, [5]uint32{item.StartPoint().Row, item.StartPoint().Column,
				item.EndPoint().Row - item.StartPoint().Row, 1, 0})
		} else if !entry.IsNull() && entry.StartPoint().Row != entry.EndPoint().Row {
			logger.Warn(`envd import "include()" not in single line"`, zap.Uint32("start", item.StartPoint().Row),
				zap.Uint32("end", item.EndPoint().Row))
		}
	}

	// sort absolute view by lineno to create relative view
	sort.Slice(absolute_view, func(i, j int) bool {
		return absolute_view[i][0] < absolute_view[j][0]
	})

	// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocument_semanticTokens
	// calculate relative flatten view
	// [{deltaLine:2, deltaStartChar:5, length:3, tokenType:0, tokenModifiers:3}] -> [2,5,3,0,3]
	relative_view := []uint32{}

	last_line := uint32(0)
	last_start := uint32(0)
	for _, item := range absolute_view {
		var deltaStartChar uint32
		deltaLine := item[0] - last_line
		if last_line != item[0] {
			deltaStartChar = item[1]
			last_line = item[0]
			last_start = 0
		} else {
			deltaStartChar = item[1] - last_start
			last_start = item[1]
		}
		relative_view = append(relative_view, deltaLine, deltaStartChar, item[2], item[3], item[4])
	}

	return relative_view
}
