package main

import (
	"context"
	"os"

	"github.com/mgjules/logem"
	"golang.org/x/exp/slog"
)

func main() {
	// Init OTEL trace provider.
	// initTraceProvider()

	// Create logger using logem.Handler.
	logger := slog.New(
		logem.NewHandler(
			slog.NewTextHandler(os.Stdout),
			logem.WithMinLevel(slog.LevelInfo),
			logem.WithStackTrace(true),
			logem.WithTraceID(true),
			logem.WithSpanID(true),
		),
	)
	slog.SetDefault(logger)

	// Use logger to log messages, etc.
	// Please ensure that the context being passed has proper trace information.
	ctx := context.TODO()
	logger.WithContext(ctx).Info("hello", "count", 3)
}
