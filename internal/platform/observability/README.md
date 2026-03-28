# Platform Observability

**Path:** `internal/platform/observability`

## Purpose

This package owns observability bootstrap helpers for `aion-api`.
It initializes trace and metric exporters, normalizes OTLP settings, and defines the shared resource metadata attached to telemetry.

## Structure

| Path | Responsibility |
| --- | --- |
| `tracer/` | OTLP HTTP trace exporter bootstrap and global tracer provider |
| `metric/` | OTLP HTTP metric exporter bootstrap and global meter provider |
| `helpers.go` | shared header parsing and endpoint normalization helpers |

## Current Runtime Behavior

- tracing and metrics both bootstrap from `cfg.Observability`
- OTLP endpoints accept `host:port` or full `http(s)://...` values and are normalized before exporter creation
- resource attributes include:
  - `service.name`
  - `service.version`
  - `deployment.environment`
  - `host.name`
  - `service.instance.id`
- tracing installs W3C `TraceContext` + `Baggage` propagators globally
- if exporter initialization fails, the app degrades gracefully and returns a no-op cleanup function instead of aborting startup

## Boundaries

- Instrumentation points live in handlers, controllers, usecases, repositories, and runtime boundaries outside this package.
- Collector/container wiring belongs to `infrastructure/observability`.
- Sampling and exporter behavior remain config-driven; do not hardcode environment-specific endpoints elsewhere.

## Related Docs

- [`../../../infrastructure/observability/README.md`](../../../infrastructure/observability/README.md)
- [`../config/README.md`](../config/README.md)

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
