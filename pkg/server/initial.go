package server

import (
	"context"

	"github.com/tensorchord/envd-lsp/pkg/query"
	"go.lsp.dev/protocol"
)

type SemanticTokensOptions struct {
	WorkDoneProgress bool                          `json:"workDoneProgress,omitempty"`
	Legend           protocol.SemanticTokensLegend `json:"legend,omitempty"`
	Range            bool                          `json:"range,omitempty"`
	Full             bool                          `json:"full,omitempty"`
}

func (s *Server) Initialize(ctx context.Context,
	params *protocol.InitializeParams) (result *protocol.InitializeResult, err error) {
	result, err = s.Server.Initialize(ctx, params)
	result.Capabilities.SemanticTokensProvider = SemanticTokensOptions{
		Legend: protocol.SemanticTokensLegend{
			TokenTypes:     query.TokenTypes,
			TokenModifiers: []protocol.SemanticTokenModifiers{},
		},
		Full:  true,
		Range: false,
	}
	return
}
