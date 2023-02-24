package logem

import "golang.org/x/exp/slog"

// Option applies optional configuration for a Handler.
type Option func(h *Handler)

// WithMinLevel sets the minimum level on which log messages are handled.
func WithMinLevel(level slog.Leveler) Option {
	return func(h *Handler) {
		h.minLevel = level
	}
}

// WithStackTrace adds a stack trace to traced log messages.
func WithStackTrace(on bool) Option {
	return func(h *Handler) {
		h.withStackTrace = on
	}
}

// WithTraceID adds the Trace ID to log messages.
func WithTraceID(on bool) Option {
	return func(h *Handler) {
		h.withTraceID = on
	}
}

// WithSpanID adds the Span ID to log messages.
func WithSpanID(on bool) Option {
	return func(h *Handler) {
		h.withSpanID = on
	}
}
