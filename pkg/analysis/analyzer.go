// Copyright 2022 The envd Authors
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

// package lsp is mainly copied from https://github.com/tilt-dev/starlark-lsp/blob/main/pkg/cli/start.go
package analysis

import (
	"context"
	"io/fs"

	starlark_analyzer "github.com/tilt-dev/starlark-lsp/pkg/analysis"
)

type AnalyzerOption = starlark_analyzer.AnalyzerOption
type BuiltinAnalyzerOptionProvider = func() starlark_analyzer.AnalyzerOption

type Analyzer struct {
	*starlark_analyzer.Analyzer
}

func NewAnalyzer(ctx context.Context, opts ...starlark_analyzer.AnalyzerOption) (*Analyzer, error) {
	analysis, err := starlark_analyzer.NewAnalyzer(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &Analyzer{analysis}, nil
}

// Reference of starlark-lsp/pkg/analysis/builtins
func WithBuiltins(f fs.FS) AnalyzerOption {
	return starlark_analyzer.WithBuiltins(f)
}

func WithStarlarkBuiltins() AnalyzerOption {
	return starlark_analyzer.WithStarlarkBuiltins()
}
