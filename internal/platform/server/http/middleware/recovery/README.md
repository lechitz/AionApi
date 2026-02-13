# HTTP Recovery Middleware

**Path:** `internal/platform/server/http/middleware/recovery`

## Overview

This package provides panic-recovery middleware for the HTTP server.
Its only responsibility is to intercept panics in the request pipeline and delegate standardized error handling to the generic recovery handler.

## Package Scope

| Area | Responsibility |
| --- | --- |
| Panic interception | Catch runtime panics from downstream handlers/middlewares |
| Error correlation | Generate a unique error identifier per recovered panic |
| Delegation | Forward recovered payload to `generic/handler.RecoveryHandler` |
| Safety | Prevent process crash and preserve API availability |

## Files

| File | Purpose |
| --- | --- |
| `recover_middleware.go` | Exposes `New(recoveryHandler)` middleware with `defer` + `recover` flow |

## Public API Reference

| Function | Returns | Description |
| --- | --- | --- |
| `New(recoveryHandler *handler.Handler)` | `func(http.Handler) http.Handler` | Wraps request execution and delegates panic handling to generic recovery handler |

## Runtime Flow

1. Request enters middleware.
2. Middleware executes downstream chain with a `defer` recovery guard.
3. If a panic occurs, middleware generates `errorID` using `uuid.New().String()`.
4. Middleware calls `recoveryHandler.RecoveryHandler(w, r, rec, errorID)`.
5. The generic handler writes the sanitized HTTP 500 response and telemetry/log details.

## Usage

```go
r.Use(recovery.New(genericRecoveryHandler))
```

## Design Notes

- Keep recovery middleware thin: detect panic and delegate; do not format responses directly here.
- Place this middleware high in the chain so panics from subsequent middlewares are also recovered.
- Keep business/domain behavior outside middleware scope.

## Package Improvements

- Add dedicated middleware tests to validate panic interception and delegation behavior.
- Validate nil `recoveryHandler` at construction time or document expected non-nil contract explicitly.
- Consider exposing a deterministic `errorID` generator interface for easier testability.
- Add a short package-level example showing recommended middleware ordering in the HTTP stack.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../../README.md)
<!-- doc-nav:end -->
