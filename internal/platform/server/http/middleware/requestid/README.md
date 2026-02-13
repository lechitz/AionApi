# HTTP Request ID Middleware

**Path:** `internal/platform/server/http/middleware/requestid`

## Overview

This package guarantees a valid request correlation ID for every HTTP request.
It normalizes `X-Request-ID`, stores it in request context, and echoes it in the response header.

## Package Scope

| Area | Responsibility |
| --- | --- |
| Header normalization | Read and sanitize incoming `X-Request-ID` |
| Fallback generation | Create UUID when incoming ID is missing or invalid |
| Context propagation | Store final ID under `ctxkeys.RequestID` |
| Response propagation | Set the same value in response header for client correlation |

## Files

| File | Purpose |
| --- | --- |
| `request_id_middleware.go` | Request ID validation/generation and context/header propagation |

## Public API Reference

| Function | Returns | Description |
| --- | --- | --- |
| `New()` | `func(http.Handler) http.Handler` | Middleware that ensures valid request ID and propagates it |
| `isUUID(s string)` | `bool` | Internal helper for UUID format validation |

## Runtime Behavior

1. Read `X-Request-ID` from request headers (`commonkeys.XRequestID`).
2. Trim whitespace.
3. If missing, too long (`>128`), or not UUID, replace with `uuid.NewString()`.
4. Inject final value into request context (`ctxkeys.RequestID`).
5. Set response header `X-Request-ID` with the same final value.
6. Continue request chain with enriched context.

## Usage

```go
r.Use(requestid.New())
```

## Accessing Request ID in Handlers

```go
rid, _ := r.Context().Value(ctxkeys.RequestID).(string)
logger.Infow("request received", "request_id", rid)
```

## Design Notes

- Keep this middleware early in the HTTP chain so every downstream log/trace can use the same request ID.
- Keep request ID semantics transport-level only; no domain coupling.
- Use shared constants (`commonkeys`, `ctxkeys`) to avoid key drift.

## Package Improvements

- Remove the unreachable `else if len(reqID) > maxLen` branch, since the previous `if` already handles `len(reqID) > maxLen`.
- Add unit tests for: valid incoming UUID, invalid UUID fallback, oversized header fallback, and whitespace-only header.
- Consider accepting additional safe ID formats if interoperability with upstream gateways becomes necessary.
- Add explicit tracing helper guidance for binding `ctxkeys.RequestID` into span attributes in adapters.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../../README.md)
<!-- doc-nav:end -->
