package observability

import (
	"context"

	"github.com/lechitz/AionApi/cmd/aion-api/constants"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/def"
	"github.com/lechitz/AionApi/internal/platform/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
)

// InitTracer initializes the OpenTelemetry tracer using the provided configuration.
// It returns a cleanup function that shuts down the tracer provider and any associated resources.
func InitTracer(cfg config.Config, logger output.Logger) func() {
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(cfg.Observability.OtelExporterOTLPEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		logger.Errorw(constants.ErrInitializeOTPL, def.Error, err)
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
			logger.Errorw(constants.ErrFailedToShutdownTracerProvider, def.Error, err)
		}
	}
}
