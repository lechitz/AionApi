# infrastructure/observability/otel

OpenTelemetry Collector configuration for AionAPI. This collector receives OTLP from the API and routes traces/metrics to the rest of the stack.

## Package Composition

- `otel-collector-config.yaml`
  - Receivers, processors, and exporters for traces and metrics.

## Flow (Where it comes from -> Where it goes)

AionAPI (OTLP) -> OTel Collector -> exporters (Prometheus, Jaeger/Zipkin)

## Why It Was Designed This Way

- Centralize telemetry routing in a single, versioned config.
- Allow multiple backends without changing the API.
- Keep local observability reproducible.

## Recommended Practices Visible Here

- Filter low-signal spans (health checks) without losing critical traces.
- Keep attribute keys aligned with `internal/shared/constants/tracingkeys`.
- Validate collector health before starting the API.

## Differentials

- Multi-exporter routing from one OTLP intake.
- Noise filtering baked into the collector pipeline.

## What Should NOT Live Here

- Application instrumentation code.
- Grafana dashboards or datasource configs.
- Secrets or production credentials.
