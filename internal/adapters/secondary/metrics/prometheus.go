// Package metrics provides functionality for exporting metrics to Prometheus.
package metrics

import (
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

// PrometheusExporter provides a Prometheus exporter instance and meter provider.
type PrometheusExporter struct {
	Exporter *prometheus.Exporter
	Provider *metric.MeterProvider
}

// NewPrometheusExporter initializes and returns a new Prometheus exporter instance.
func NewPrometheusExporter() (*PrometheusExporter, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
	)

	return &PrometheusExporter{
		Exporter: exporter,
		Provider: provider,
	}, nil
}
