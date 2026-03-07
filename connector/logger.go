package connector

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/sdk/log"
)

type loggerContext struct{}

var loggerContextKey = loggerContext{}

// LogHandler wraps slog logger with the OpenTelemetry logs exporter handler.
type LogHandler struct {
	otelHandler slog.Handler
	stdHandler  slog.Handler
}

func createLogHandler(
	serviceName string,
	logger *slog.Logger,
	provider *log.LoggerProvider,
) slog.Handler {
	options := []otelslog.Option{}
	if provider != nil {
		options = append(options, otelslog.WithLoggerProvider(provider))
	}

	otelHandler := otelslog.NewHandler(serviceName, options...)
	loggerHandler := logger.Handler()

	return LogHandler{
		otelHandler: otelHandler,
		stdHandler:  loggerHandler,
	}
}

// Enabled reports whether the handler handles records at the given level.
func (l LogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return l.stdHandler.Enabled(ctx, level)
}

// Handle handles the Record.
// It will only be called when Enabled returns true.
func (l LogHandler) Handle(ctx context.Context, record slog.Record) error {
	_ = l.stdHandler.Handle(ctx, record)

	return l.otelHandler.Handle(ctx, record)
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
// The Handler owns the slice: it may retain, modify or discard it.
func (l LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return LogHandler{
		otelHandler: l.otelHandler.WithAttrs(attrs),
		stdHandler:  l.stdHandler.WithAttrs(attrs),
	}
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
func (l LogHandler) WithGroup(name string) slog.Handler {
	return LogHandler{
		otelHandler: l.otelHandler.WithGroup(name),
		stdHandler:  l.stdHandler.WithGroup(name),
	}
}

// GetLogger gets the logger instance from context.
func GetLogger(ctx context.Context) *slog.Logger {
	value := ctx.Value(loggerContextKey)
	if value != nil {
		if logger, ok := value.(*slog.Logger); ok {
			return logger
		}
	}

	return slog.New(createLogHandler("hasura-ndc-go", slog.Default(), nil))
}

// NewContextLogger creates a new context with a logger set.
func NewContextLogger(parentContext context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(parentContext, loggerContextKey, logger)
}
