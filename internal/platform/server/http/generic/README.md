# Platform HTTP Generic Layer

**Path:** `internal/platform/server/http/generic`

## Overview

This package contains platform-level generic HTTP components shared by all contexts.
It provides health and fallback handlers (`404`, `405`, error, panic recovery) with consistent response envelopes, logging, and tracing.

## Package Scope

| Area | Responsibility |
| --- | --- |
| Health endpoint | Return service status and runtime metadata |
| Router fallback handlers | Standardize `Not Found`, `Method Not Allowed`, and generic error behavior |
| Panic recovery integration | Convert recovered panics into safe HTTP error responses |
| Shared transport DTOs | Define transport-only payloads used by generic handlers |

## Subpackages

| Subpackage | Role |
| --- | --- |
| `dto/` | Generic HTTP DTOs (currently health response payload) |
| `handler/` | Generic handler implementation and tracing/logging constants |

## Main Components

| Component | Description |
| --- | --- |
| `handler.New(logger, generalCfg)` | Creates `*handler.Handler` with logger and general app metadata |
| `(*Handler).HealthCheck` | Handles `/health` and returns app metadata + healthy status |
| `(*Handler).NotFoundHandler` | Standardized JSON `404` |
| `(*Handler).MethodNotAllowedHandler` | Standardized JSON `405` |
| `(*Handler).ErrorHandler` | Standardized JSON `500` for router-level errors |
| `(*Handler).RecoveryHandler` | Handles panic recovery payloads and emits telemetry |

## Integration Flow

1. HTTP composer creates `genericHandler := handler.New(log, cfg.General)`.
2. Router wiring connects:
- `/health` to `genericHandler.HealthCheck`
- `SetNotFound` to `genericHandler.NotFoundHandler`
- `SetMethodNotAllowed` to `genericHandler.MethodNotAllowedHandler`
- `SetError` to `genericHandler.ErrorHandler`
3. Recovery middleware delegates panic handling to `genericHandler.RecoveryHandler`.

## Observability Behavior

| Aspect | Behavior |
| --- | --- |
| Tracing | Creates spans for health, error, and recovery flows |
| Span attributes | Includes request metadata (`path`, `request_id`, `ip`, `user_agent`) |
| Logging | Emits structured logs with shared keys from `commonkeys`/`tracingkeys` |
| Responses | Uses `httpresponse` helpers for consistent JSON envelopes |

## Design Notes

- This layer is transport/platform only and should not contain business rules.
- Generic handlers provide a stable baseline for all bounded contexts.
- Constants in `handler/0_generic_handler_constants.go` centralize trace/log message semantics.

## Package Improvements

- Add dedicated tests for each generic handler (`health`, `404`, `405`, `error`, `recovery`) to lock response and tracing contracts.
- Replace duplicated manual error body construction in `NotFoundHandler` and `MethodNotAllowedHandler` with a shared `httpresponse.WriteError` flow for consistency.
- Consider moving tracer names and common messages to `internal/shared/constants/tracingkeys` when broadly reused outside this package.
- Evaluate whether `HealthCheck` should support `HEAD` with empty body explicitly, depending on monitoring tool expectations.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../README.md)
<!-- doc-nav:end -->
