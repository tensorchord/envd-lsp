# Language server for envd

A language server for envd, a Python-inspired configuration language.

envd-lsp uses [starlark-lsp](https://github.com/tilt-dev/starlark-lsp/), [go.lsp.dev](https://go.lsp.dev/) and [Tree sitter](https://tree-sitter.github.io/tree-sitter/) as its main dependencies to implement the LSP/JSON-RPC protocol and envd language analysis, respectively.

## Build from source

```
make
```

## CLI

```
NAME:
   envd-lsp - language server for envd

USAGE:
   envd-lsp [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --address value  Address (hostname:port) to listen on
   --debug          Enable debug logging (default: false)
   --help, -h       show help (default: false)
   --version, -v    print the version (default: false)
```
