package server

import (
	"context"
	"fmt"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

func (s Server) SemanticTokensFull(ctx context.Context, params *protocol.SemanticTokensParams) (result *protocol.SemanticTokens, err error) {
	doc, err := s.docs.Read(ctx, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}
	defer doc.Close()

	logger := protocol.LoggerFromContext(ctx)
	logger.Debug("semantic token", zap.String("path", string(doc.URI())))

	data := s.analyzer.SemanticToken(ctx, doc)
	logger.With(zap.Namespace("semantic token")).Debug(fmt.Sprintf("found semantic tokens: %v", data))

	return &protocol.SemanticTokens{
		ResultID: "",
		Data:     data,
	}, nil
}
