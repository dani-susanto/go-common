package log

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/trace"
)

type Log interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type log struct {
	logger *slog.Logger
}

type otelHandler struct {
	slog.Handler
}

func (h *otelHandler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		r.AddAttrs(
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}
	return h.Handler.Handle(ctx, r)
}

func New(serviceName string) Log {
	logger := slog.New(slog.NewMultiHandler(
		&otelHandler{
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		},
		otelslog.NewHandler(serviceName),
	))

	return &log{logger: logger}
}

func (l *log) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *log) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *log) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *log) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}
