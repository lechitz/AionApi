# Instrumented HTTP Client

**Path:** `internal/platform/httpclient`

## Overview

Shared outbound HTTP client factory with OpenTelemetry instrumentation.
Used by secondary adapters that call external HTTP services.

## Responsibilities

| Area | Responsibility |
| --- | --- |
| Outbound tracing | Attach spans and propagate trace context |
| Client standardization | Centralize timeout/transport defaults |
| DI integration | Provide reusable client via Fx wiring |

## Usage Pattern

- Inject `*http.Client` into adapters.
- Keep service-specific logic in adapter packages.

## Design Notes

- Keep this package generic and transport-level.
- Avoid embedding service-specific URLs/protocol logic.
- Prefer constructor injection over ad-hoc client creation.

## Package Improvements

- Add test helpers for deterministic transport mocking.
- Add policy for default timeout values by environment.
- Add sample instrumentation verification test.
- Add guidance for retries/circuit breaking integration.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
