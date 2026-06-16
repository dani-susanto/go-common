package log

import (
	"context"
	"fmt"
	golog "log"
	"strings"

	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
)

const (
	red   = "\033[31m"
	green = "\033[32m"
	reset = "\033[0m"
)

type Log interface {
	Info(v ...any)
	Error(v ...any)
	Infof(format string, v ...any)
	Errorf(format string, v ...any)
}

type log struct {
	ctx        context.Context
	otelLogger otellog.Logger
}

func New(ctx context.Context, serviceName string) Log {
	return &log{
		ctx:        ctx,
		otelLogger: global.GetLoggerProvider().Logger(serviceName),
	}
}

func (l *log) Info(v ...any) {
	l.print(otellog.SeverityInfo, green, "[INF]", sprint(v))
}

func (l *log) Error(v ...any) {
	l.print(otellog.SeverityError, red, "[ERR]", sprint(v))
}

func (l *log) Infof(format string, v ...any) {
	l.print(otellog.SeverityInfo, green, "[INF]", fmt.Sprintf(format, v...))
}

func (l *log) Errorf(format string, v ...any) {
	l.print(otellog.SeverityError, red, "[ERR]", fmt.Sprintf(format, v...))
}

func (l *log) print(severity otellog.Severity, color, prefix, msg string) {
	golog.Print(color + prefix + reset + " " + msg)
	l.emit(severity, msg)
}

func (l *log) emit(severity otellog.Severity, msg string) {
	var r otellog.Record
	r.SetSeverity(severity)
	r.SetBody(otellog.StringValue(msg))
	l.otelLogger.Emit(l.ctx, r)
}

func sprint(v []any) string {
	return strings.TrimSuffix(fmt.Sprintln(v...), "\n")
}
