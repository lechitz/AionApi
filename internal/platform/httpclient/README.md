# Instrumented HTTP Client

**Path:** `internal/platform/httpclient`

## Purpose

This package builds the shared outbound HTTP client used by secondary adapters.

## Current Flow

| Piece | Responsibility |
| --- | --- |
| `NewInstrumentedClient` | wrap the base transport with `otelhttp` unless instrumentation is disabled |
| `NewClient` | expose the stdlib client through the `platform/ports/output/httpclient.HTTPClient` interface |
| `fxapp.InfraModule` | provide the client with timeout derived from config (`AionChat.Timeout` fallback) |

## Boundaries

- adapters should depend on the output port, not on raw `*http.Client`
- service-specific URLs, payload semantics, and retries belong in the owning adapter
- this package owns transport instrumentation and generic request helpers only

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
