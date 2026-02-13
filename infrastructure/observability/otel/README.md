# OpenTelemetry Collector Config

**Path:** `infrastructure/observability/otel`

## Overview

Collector pipeline configuration for traces and metrics emitted by Aion services.
It centralizes telemetry routing to downstream backends.

## Files

| File | Purpose |
| --- | --- |
| `otel-collector-config.yaml` | Receivers, processors, exporters for telemetry pipeline |

## Design Notes

- Keep exporter endpoints configurable per environment.
- Keep attribute naming aligned with `tracingkeys` conventions.
- Tune processors to reduce noise while preserving signal.

## Package Improvements

- Add pipeline diagram (receiver -> processor -> exporter).
- Add config validation checks in CI.
- Add troubleshooting matrix for exporter connectivity issues.
- Add documented defaults for sampling/timeout behavior.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
