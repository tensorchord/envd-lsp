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

package main

import (
	"fmt"
	"os"

	cli "github.com/urfave/cli/v2"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/tensorchord/envd-lsp/pkg/api"
	"github.com/tensorchord/envd-lsp/pkg/lsp"
	"github.com/tensorchord/envd-lsp/pkg/version"
)

// go lsp server uses zap, thus we have to keep two loggers (logrus and zap) in envd.
// zap is only used in lsp subcommand.
var logLevel = zap.NewAtomicLevelAt(zapcore.WarnLevel)

func startLSP(clicontext *cli.Context) error {
	logger, cleanup := newzapLogger()
	defer cleanup()
	ctx := protocol.WithLogger(clicontext.Context, logger)

	apiVersion := api.APIOptions[clicontext.String("api")]
	s := lsp.New(apiVersion)
	err := s.Start(ctx, clicontext.String("address"))
	return err
}

func newzapLogger() (logger *zap.Logger, cleanup func()) {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = logLevel
	cfg.Development = false
	logger, err := cfg.Build()
	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %v", err))
	}

	cleanup = func() {
		_ = logger.Sync()
	}
	return logger, cleanup
}

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(c.App.Name, version.Package, version.GetVersion().String())
	}

	app := cli.NewApp()
	app.Name = "envd-lsp"
	app.Usage = "language server for envd"
	app.Version = version.GetVersion().String()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Usage: "Address (hostname:port) to listen on",
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug logging",
		},
		&cli.StringFlag{
			Name:   "api",
			Usage:  "API version to use of envd",
			Value:  "stable",
			Action: api.ArgValidator,
		},
	}

	// Deal with debug flag.
	var debugEnabled bool

	app.Before = func(context *cli.Context) error {
		debugEnabled = context.Bool("debug")

		if debugEnabled {
			logLevel.SetLevel(zapcore.DebugLevel)
		}
		return nil
	}

	app.Action = startLSP
	handleErr(debugEnabled, app.Run(os.Args))
}

func handleErr(debug bool, err error) {
	if err == nil {
		return
	}
	if debug {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	os.Exit(1)
}
