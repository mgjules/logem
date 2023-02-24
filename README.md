# Logem

[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://godoc.org/github.com/mgjules/logem)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=for-the-badge)](LICENSE)

Logem, short for "Let's Log em!", is an opinionated handler wrapper for [slog](https://pkg.go.dev/golang.org/x/exp/slog).

It combines a LevelHandler and an TraceHandler in a single, simple to use handler.

## Disclaimer

This package is not suited for production use right now. It is heavily being worked on.

However feel free to inspire yourself from it or fork it.

## Installation

```shell
go install github.com/mgjules/logem@v0.0.1
```

## Usage

```go
package main

import (
  "context"
  "os"

  "github.com/mgjules/logem"
  "golang.org/x/exp/slog"
)

func main() {
  // Init OTEL tracer...
  // initTracer()

  // Create logger using logem.Handler.
  logger := slog.New(logem.NewHandler(slog.NewTextHandler(os.Stdout)))
  slog.SetDefault(logger)

  // Use logger to log messages, etc.
  // Please ensure that context being passed has proper trace information.
  ctx := context.TODO()
  logger.WithContext(ctx).Info("hello", "count", 3)
}
```

## Stability

This project follows [SemVer](http://semver.org/) strictly and is not yet `v1`.

Breaking changes might be introduced until `v1` is released.

This project follows the [Go Release Policy](https://golang.org/doc/devel/release.html#policy). Each major version of Go is supported until there are two newer major releases.