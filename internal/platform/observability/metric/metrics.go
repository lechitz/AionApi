// Package metric provides utilities for initializing and managing OpenTelemetry metrics.
package metric

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/observability"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
)

const (
	// ErrFailedToInitializeOTLPMetricsExporter is logged when the OTLP metrics exporter cannot be created.
	ErrFailedToInitializeOTLPMetricsExporter = "failed to initialize OTLP metric exporter"

	// ErrInvalidOTELExporterTimeout is logged when the timeout string cannot be parsed as a valid duration.
	ErrInvalidOTELExporterTimeout = "invalid OTLP exporter timeout"

	// ErrToShutDownOTELMetrics is logged when shutting down the metrics provider fails.
	ErrToShutDownOTELMetrics = "failed to shutdown OTLP metrics provider"
)

// InitOtelMetrics sets up the OpenTelemetry MeterProvider using the given configuration,
// installs it as the global provider, and returns a cleanup function to gracefully shut it down.
func InitOtelMetrics(cfg *config.Config, logger logger.ContextLogger) func() {
	opts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(cfg.Observability.OtelExporterOTLPEndpoint),
	}

	if cfg.Observability.OtelExporterInsecure {
		opts = append(opts, otlpmetrichttp.WithInsecure())
	}

	if cfg.Observability.OtelExporterTimeout != "" {
		timeout, err := time.ParseDuration(cfg.Observability.OtelExporterTimeout)
		if err == nil {
			opts = append(opts, otlpmetrichttp.WithTimeout(timeout))
		} else {
			logger.Warnw(ErrInvalidOTELExporterTimeout, commonkeys.Error, err)
		}
	}

	if cfg.Observability.OtelExporterCompression == "gzip" {
		opts = append(opts, otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression))
	}

	headers := observability.ParseHeaders(cfg.Observability.OtelExporterHeaders)
	if len(headers) > 0 {
		opts = append(opts, otlpmetrichttp.WithHeaders(headers))
	}

	exporter, err := otlpmetrichttp.New(context.Background(), opts...)
	if err != nil {
		logger.Errorw(ErrFailedToInitializeOTLPMetricsExporter, commonkeys.Error, err)
		panic(err)
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.Observability.OtelServiceName),
			semconv.ServiceVersionKey.String(cfg.Observability.OtelServiceVersion),
		)),
	)

	otel.SetMeterProvider(provider)

	return func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			logger.Errorw(ErrToShutDownOTELMetrics, commonkeys.Error, err)
		}
	}
}
