// Package tracer provides a wrapper for the OpenTelemetry tracer.
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
	// ErrInitializeOTPL is the error message for when the OTLP exporter fails to initialize.
	ErrInitializeOTPL = "failed to initialize OTLP exporter"

	// ErrFailedToShutdownTracerProvider is the error message for when the tracer token fails to shutdown.
	ErrFailedToShutdownTracerProvider = "failed to shutdown tracer token"
)

// InitTracer initializes the OpenTelemetry tracer using the provided configuration.
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
			logger.Warnw("Invalid OTEL exporter timeout, using default", commonkeys.Error, err)
		}
	}

	if cfg.Observability.OtelExporterCompression == "gzip" {
		opts = append(opts, otlptracehttp.WithCompression(otlptracehttp.GzipCompression))
	}

	headers := observability.ParseHeaders(cfg.Observability.OtelExporterHeaders)
	if len(headers) > 0 {
		opts = append(opts, otlptracehttp.WithHeaders(headers))
	}

	exporter, err := otlptracehttp.New(context.Background(), opts...)
	if err != nil {
		logger.Errorw(ErrInitializeOTPL, commonkeys.Error, err)
		panic(err) // TODO: avaliar se mantem panic ou substitui por os.Exit
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
