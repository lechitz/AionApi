# infrastructure/observability

Observability infrastructure for AionAPI. This folder keeps the logging, metrics, and tracing stack reproducible across environments.

## Package Composition

- `otel/`
  - OTel Collector routing for traces and metrics.
- `prometheus/`
  - Scrape configuration for metrics.
- `grafana/`
  - Datasources, dashboards, and provisioning.
- `loki/`
  - Log storage configuration.
- `fluentbit/`
  - Log collection and routing.
- `scripts/`
  - Automation helpers for local setup and validation.

## Flow (Where it comes from -> Where it goes)

AionAPI telemetry -> OTel/Fluent Bit -> Prometheus/Loki/Jaeger -> Grafana

## Why It Was Designed This Way

- Keep observability assets versioned with the codebase.
- Make local setup deterministic and fast.
- Separate collection, storage, and visualization concerns.

## Recommended Practices Visible Here

- Align metric and trace names with `internal/platform/observability`.
- Keep datasource URLs and ports environment-specific.
- Validate endpoints after configuration changes.

## Differentials

- Full stack provisioning (collector, metrics, logs, dashboards) in-repo.
- Automation scripts for quick validation.

## What Should NOT Live Here

- Application instrumentation code.
- Business logic or runtime configuration.
- Secrets or production credentials.
