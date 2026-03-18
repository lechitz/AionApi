# OpenTelemetry Collector Config

**Path:** `infrastructure/observability/otel`

## Purpose

This package routes OTLP telemetry emitted by Aion services.

## Current Pipeline

| Stage | Behavior |
| --- | --- |
| Receiver | accepts OTLP over HTTP and gRPC |
| Processor | filters health spans and batches telemetry |
| Exporter | sends traces to Jaeger and exposes metrics for Prometheus scraping |

## Source File

| File | Responsibility |
| --- | --- |
| `otel-collector-config.yaml` | receiver, processor, and exporter wiring |

## Boundaries

- attribute naming and span semantics remain owned by the application code
- collector routing lives here, not in Grafana or Prometheus
- if exporter endpoints change in compose, this config must change in the same PR

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
