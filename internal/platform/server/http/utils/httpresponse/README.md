# HTTP Response Utilities

**Path:** `internal/platform/server/http/utils/httpresponse`

## Overview

This package standardizes HTTP JSON responses across handlers, including success envelopes, error envelopes, and span-aware error helpers.
It is the response boundary for HTTP adapters and relies on semantic errors from `sharederrors`.

## Package Scope

| Area | Responsibility |
| --- | --- |
| JSON writing | Encode payloads and set HTTP headers consistently |
| Success responses | Build normalized success envelope (`ResponseBody`) |
| Error responses | Map semantic errors to status and write normalized error envelope |
| Span-aware helpers | Record OTel span status/attributes before returning HTTP errors |

## Files

| File | Purpose |
| --- | --- |
| `httpresponse.go` | Response envelope, writers, and span-aware helpers |
| `httpresponse_test.go` | Table-driven tests for response/status/headers and tracing behavior |

## Response Envelope

| Field | Type | Description |
| --- | --- | --- |
| `date` | `time.Time` | UTC timestamp for response generation |
| `result` | `any` | Success payload |
| `message` | `string` | Optional success message |
| `error` | `string` | Client-facing error message |
| `details` | `string` | Raw error detail (`err.Error()`) |
| `code` | `int` | HTTP status code |

## Public API Reference

### Core Writers

| Function | Behavior |
| --- | --- |
| `WriteJSON(w, status, payload, headers...)` | Writes raw JSON payload; skips body for `204 No Content` |
| `WriteSuccess(w, status, result, message, headers...)` | Writes standardized success envelope |
| `WriteError(w, err, message, log, headers...)` | Maps error status, logs (`Errorw`), writes standardized error envelope |
| `WriteDecodeError(w, err, log, headers...)` | Shortcut for malformed body (`Invalid request body`) |
| `WriteAuthError(w, err, log, headers...)` | Shortcut for auth failures (`Unauthorized`) |
| `WriteNoContent(w, headers...)` | Writes `204` with optional headers |

### Span-aware Writers

| Function | Trace behavior |
| --- | --- |
| `WriteAuthErrorSpan(...)` | Records error, sets span status error, sets `http.status_code`, then writes auth error |
| `WriteDecodeErrorSpan(...)` | Records decode error and returns `400` response |
| `WriteValidationErrorSpan(...)` | Records validation error and returns error response using `err.Error()` as message |
| `WriteDomainErrorSpan(...)` | Records domain error and maps status via `sharederrors.MapErrorToHTTPStatus` |

## Status Mapping Behavior

`WriteError` and domain span helpers delegate status mapping to:

- `sharederrors.MapErrorToHTTPStatus(err)`

This keeps status semantics centralized and consistent across all HTTP adapters.

## Tested Behaviors

| Behavior | Verified by |
| --- | --- |
| `204` responses do not write body/content-type | `TestWriteJSON`, `TestWriteNoContent` |
| Success envelope has expected code/message/result/date | `TestWriteSuccess` |
| Error envelope contains `error`, `details`, and mapped status | `TestWriteError`, `TestWriteAuthAndDecodeError` |
| Custom response headers are preserved | `TestWriteError_WithCustomHeaders` |
| Span helpers set `codes.Error` and `tracingkeys.HTTPStatusCodeKey` | `TestSpanErrorResponses` |

## Usage Example

```go
if err != nil {
    httpresponse.WriteError(w, err, "failed to create resource", log)
    return
}

httpresponse.WriteSuccess(w, http.StatusCreated, result, "resource created")
```

## Design Notes

- Keep handlers thin by delegating response shape and status mapping to this package.
- Keep transport concerns here; domain usecases should not depend on this package.
- Use `Write*Span` helpers at adapter boundaries where tracing is active.

## Package Improvements

- Evaluate if `details` should be conditionally redacted in production paths to reduce sensitive error leakage risk.
- Consider adding a helper for `202 Accepted` async flows to keep envelope semantics explicit.
- Add explicit tests for `WriteSuccess`/`WriteError` with multiple header maps to validate merge precedence.
- Consider documenting a strict policy for client-facing `message` vs internal `details` to improve consistency.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../../README.md)
<!-- doc-nav:end -->
