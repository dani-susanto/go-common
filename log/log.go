package log

import (
	"context"
	"log/slog"
	"os"
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

func New() Log {
	return &log{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
	}
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
