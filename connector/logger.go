package connector

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hasura/ndc-sdk-go/utils"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

// LogHandler wraps slog logger with the OpenTelemetry logs exporter handler.
type LogHandler struct {
	otelHandler slog.Handler
	stdHandler  slog.Handler
}

func createLogHandler(serviceName string, logger *slog.Logger, provider *log.LoggerProvider) slog.Handler {
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

// create OpenTelemetry logger provider.
func newLoggerProvider(ctx context.Context, config *OTLPConfig, otelDisabled bool, res *resource.Resource) (*log.LoggerProvider, error) {
	logsEndpoint := utils.GetDefault(config.OtlpLogsEndpoint, config.OtlpEndpoint)
	if otelDisabled || config.LogsExporter != "otlp" || logsEndpoint == "" {
		return log.NewLoggerProvider(), nil
	}

	endpoint, protocol, insecure, err := parseOTLPEndpoint(
		logsEndpoint,
		utils.GetDefault(config.OtlpLogsProtocol, config.OtlpProtocol),
		utils.GetDefaultPtr(config.OtlpLogsInsecure, config.OtlpInsecure),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OTLP logs endpoint: %w", err)
	}

	compressorStr, compressorInt, err := parseOTLPCompression(utils.GetDefault(config.OtlpLogsCompression, config.OtlpCompression))
	if err != nil {
		return nil, fmt.Errorf("failed to parse OTLP logs compression: %w", err)
	}

	opts := []log.LoggerProviderOption{log.WithResource(res)}

	if protocol == otlpProtocolGRPC {
		options := []otlploggrpc.Option{
			otlploggrpc.WithEndpoint(endpoint),
			otlploggrpc.WithCompressor(compressorStr),
		}

		if insecure {
			options = append(options, otlploggrpc.WithInsecure())
		}

		logExporter, err := otlploggrpc.New(ctx, options...)
		if err != nil {
			return nil, err
		}

		opts = append(opts, log.WithProcessor(log.NewBatchProcessor(logExporter)))
	} else {
		options := []otlploghttp.Option{
			otlploghttp.WithEndpoint(endpoint),
			otlploghttp.WithCompression(otlploghttp.Compression(compressorInt)),
		}
		if insecure {
			options = append(options, otlploghttp.WithInsecure())
		}

		logExporter, err := otlploghttp.New(ctx, options...)
		if err != nil {
			return nil, err
		}

		opts = append(opts, log.WithProcessor(log.NewBatchProcessor(logExporter)))
	}

	loggerProvider := log.NewLoggerProvider(opts...)

	return loggerProvider, nil
}

// GetLogger gets the logger instance from context.
func GetLogger(ctx context.Context) *slog.Logger {
	value := ctx.Value(logContextKey)
	if value != nil {
		if logger, ok := value.(*slog.Logger); ok {
			return logger
		}
	}

	return slog.New(createLogHandler("hasura-ndc-go", slog.Default(), nil))
}
