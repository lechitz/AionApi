package observability

import (
	"context"

	"github.com/lechitz/AionApi/cmd/aion-api/constants"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/def"
	"github.com/lechitz/AionApi/internal/platform/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
)

// InitOtelMetrics initializes the OpenTelemetry metrics provider using the given configuration.
// It returns a cleanup function to gracefully shut down the provider.
func InitOtelMetrics(cfg config.Config, logger output.Logger) func() {
	exporter, err := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithEndpoint(cfg.Observability.OtelExporterOTLPEndpoint),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		logger.Errorw(constants.ErrFailedToInitializeOTLPMetricsExporter, def.Error, err)
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
		_ = provider.Shutdown(context.Background())
	}
}
