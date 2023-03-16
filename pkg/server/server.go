// Copyright 2022 The envd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"

	"github.com/tensorchord/envd-lsp/pkg/analysis"
	"github.com/tilt-dev/starlark-lsp/pkg/document"
	"github.com/tilt-dev/starlark-lsp/pkg/middleware"
	starlark_server "github.com/tilt-dev/starlark-lsp/pkg/server"
	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
)

type Server struct {
	*starlark_server.Server
	// reference of starlark-lsp/pkg/document to overload private attribute
	docs *document.Manager
	// Envd overwrite analyzer
	analyzer *analysis.Analyzer
}

func NewServer(cancel context.CancelFunc, notifier protocol.Client, docManager *document.Manager, analyzer *analysis.Analyzer) *Server {
	return &Server{
		Server:   starlark_server.NewServer(cancel, notifier, docManager, analyzer.Analyzer),
		docs:     docManager,
		analyzer: analyzer,
	}
}

// Overwrite starlark-lsp handler
func (s *Server) Handler(middlewares ...middleware.Middleware) jsonrpc2.Handler {
	serverHandler := protocol.ServerHandler(s, jsonrpc2.MethodNotFoundHandler)
	return middleware.WrapHandler(serverHandler, middlewares...)
}
