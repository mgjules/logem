package logem

import (
	"context"
	"fmt"
	"runtime"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.16.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
)

// Handler wraps an slog.Handler to add trace support.
type Handler struct {
	handler slog.Handler

	minLevel       slog.Leveler
	withStackTrace bool
	withTraceID    bool
	withSpanID     bool

	opts []Option
}

// NewHandler returns a TraceHandler.
func NewHandler(h slog.Handler, opts ...Option) *Handler {
	// Optimization: avoid surface-level chains of Handler.
	if v, ok := h.(*Handler); ok {
		h = v.Handler()
	}

	handler := &Handler{
		handler: h,
		opts:    opts,
	}

	for _, opt := range opts {
		opt(handler)
	}

	return handler
}

// Enabled is a simple wrapper method around h.handler.Enabled.
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	if h.minLevel != nil {
		return level >= h.minLevel.Level()
	}

	return h.handler.Enabled(ctx, level)
}

// Handle is a simple wrapper method around h.handler.Handle.
func (h *Handler) Handle(r slog.Record) error {
	if span := trace.SpanFromContext(r.Context); span.IsRecording() {
		attrs := make([]attribute.KeyValue, 0, r.NumAttrs()+6)
		r.Attrs(func(a slog.Attr) {
			attrs = appendAttr(a, attrs, "")
		})

		attrs = append(attrs, attribute.Key("log.severity").String(r.Level.String()))
		attrs = append(attrs, attribute.Key("log.message").String(r.Message))

		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.Function != "" {
			attrs = append(attrs, semconv.CodeFunctionKey.String(f.Function))
		}
		if f.File != "" {
			attrs = append(attrs, semconv.CodeFilepathKey.String(f.File))
			attrs = append(attrs, semconv.CodeLineNumberKey.Int(f.Line))
		}

		if h.withStackTrace {
			stackTrace := make([]byte, 2048)
			n := runtime.Stack(stackTrace, false)
			attrs = append(attrs, semconv.ExceptionStacktraceKey.String(string(stackTrace[0:n])))
		}

		span.AddEvent("log", trace.WithAttributes(attrs...))

		if r.Level >= slog.LevelError {
			span.SetStatus(codes.Error, r.Message)
		}

		if h.withTraceID {
			r.AddAttrs(slog.String("trace_id", span.SpanContext().TraceID().String()))
		}
		if h.withSpanID {
			r.AddAttrs(slog.String("span_id", span.SpanContext().SpanID().String()))
		}
	}

	return h.handler.Handle(r)
}

// WithAttrs is a simple wrapper method around h.handler.WithAttrs.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewHandler(h.handler.WithAttrs(attrs), h.opts...)
}

// WithGroup is a simple wrapper method around h.handler.WithGroup.
func (h *Handler) WithGroup(name string) slog.Handler {
	return NewHandler(h.handler.WithGroup(name), h.opts...)
}

// Handler returns the Handler wrapped by h.
func (h *Handler) Handler() slog.Handler {
	return h.handler
}

func appendAttr(a slog.Attr, attrs []attribute.KeyValue, prefix string) []attribute.KeyValue {
	if a.Key == "" {
		return attrs
	}

	key := a.Key
	if prefix != "" {
		key = prefix + "." + key
	}

	value := a.Value.Resolve()

	if value.Kind() == slog.KindGroup {
		groupAttrs := value.Group()
		if len(attrs) == 0 {
			return attrs
		}

		for _, ga := range groupAttrs {
			attrs = appendAttr(ga, attrs, key)
		}

		return attrs
	}

	return append(attrs, parseSlogKeyValue(key, value))
}

func parseSlogKeyValue(key string, value slog.Value) attribute.KeyValue {
	switch value.Kind() {
	case slog.KindString:
		return attribute.String(key, value.String())
	case slog.KindInt64:
		return attribute.Int64(key, value.Int64())
	case slog.KindUint64:
		return attribute.Int64(key, int64(value.Uint64()))
	case slog.KindFloat64:
		return attribute.Float64(key, value.Float64())
	case slog.KindBool:
		return attribute.Bool(key, value.Bool())
	case slog.KindTime:
		return attribute.Int64(key, value.Time().UnixNano())
	case slog.KindDuration:
		return attribute.Int64(key, int64(value.Duration()))
	case slog.KindAny:
		return attribute.String(key, fmt.Sprint(value.Any()))
	default:
		return attribute.String(key+"_error", fmt.Sprintf("logem: unsupported value kind: %s", value.Kind()))
	}
}
