# Platform Observability

**Path:** `internal/platform/observability`

## Overview

Platform-level observability bootstrap for tracing and metrics.
This package configures telemetry providers and shared resource metadata.

## Responsibilities

| Area | Responsibility |
| --- | --- |
| Tracing setup | Initialize tracer provider/exporters |
| Metrics setup | Initialize meter provider/exporters |
| Shared metadata | Apply service/env/version resource attributes |
| Collector integration | Route telemetry to OTEL pipeline |

## Design Notes

- Keep bootstrap concerns separate from instrumentation points.
- Align telemetry keys with shared constants.
- Keep exporter configuration environment-driven.

## Package Improvements

- Add provider initialization sequence diagram.
- Add health check guidance for collector/exporters.
- Add clear defaults/sampling policy documentation.
- Add integration tests with local collector fixture.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
