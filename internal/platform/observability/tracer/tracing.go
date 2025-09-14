// Package tracer provides a wrapper for configuring and managing the OpenTelemetry tracer.
package tracer

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/observability"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
)

const (
	// ErrFailedToInitializeOTLPExporter is logged when the OTLP trace exporter cannot be created.
	ErrFailedToInitializeOTLPExporter = "failed to initialize OTLP trace exporter"

	// ErrFailedToShutdownTracerProvider is logged when the tracer provider fails to shut down.
	ErrFailedToShutdownTracerProvider = "failed to shutdown tracer provider"

	// WarnInvalidOTLPExporterTimeout is logged when the timeout string cannot be parsed and the default is used.
	WarnInvalidOTLPExporterTimeout = "invalid OTLP exporter timeout, using default"

	// CompressionGzip is the string value that enables gzip compression for the OTLP HTTP exporter.
	CompressionGzip = "gzip"
)

// InitTracer initializes the OpenTelemetry tracer using the provided configuration,
// installs it as the global tracer provider, and returns a cleanup function to shut it down gracefully.
func InitTracer(cfg *config.Config, logger logger.ContextLogger) func() {
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(cfg.Observability.OtelExporterOTLPEndpoint),
	}

	if cfg.Observability.OtelExporterInsecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	if cfg.Observability.OtelExporterTimeout != "" {
		timeout, err := time.ParseDuration(cfg.Observability.OtelExporterTimeout)
		if err == nil {
			opts = append(opts, otlptracehttp.WithTimeout(timeout))
		} else {
			logger.Warnw(WarnInvalidOTLPExporterTimeout, commonkeys.Error, err)
		}
	}

	if cfg.Observability.OtelExporterCompression == CompressionGzip {
		opts = append(opts, otlptracehttp.WithCompression(otlptracehttp.GzipCompression))
	}

	headers := observability.ParseHeaders(cfg.Observability.OtelExporterHeaders)
	if len(headers) > 0 {
		opts = append(opts, otlptracehttp.WithHeaders(headers))
	}

	exporter, err := otlptracehttp.New(context.Background(), opts...)
	if err != nil {
		logger.Errorw(ErrFailedToInitializeOTLPExporter, commonkeys.Error, err)
		panic(err)
	}

	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(cfg.Observability.OtelServiceName),
		semconv.ServiceVersionKey.String(cfg.Observability.OtelServiceVersion),
	)

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resources),
	)

	otel.SetTracerProvider(traceProvider)

	return func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			logger.Errorw(ErrFailedToShutdownTracerProvider, commonkeys.Error, err)
		}
	}
}
