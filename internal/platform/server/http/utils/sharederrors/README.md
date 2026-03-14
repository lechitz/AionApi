# HTTP Shared Errors

**Path:** `internal/platform/server/http/utils/sharederrors`

## Overview

This package centralizes HTTP-facing semantic errors and status-code mapping for handlers and adapter boundaries.
It keeps transport error behavior consistent without leaking infrastructure details.

## Package Scope

| Area | Responsibility |
| --- | --- |
| Typed errors | Define explicit error types (`ValidationError`, `UnauthorizedError`, etc.) |
| Sentinel errors | Provide reusable `errors.Is` targets for common conflicts/validation failures |
| Status mapping | Convert known errors to stable HTTP status codes |
| Message constants | Keep shared error messages consistent and reusable |

## Files

| File | Purpose |
| --- | --- |
| `errors.go` | Error constants, custom error types, constructors, and sentinel errors |
| `map_error_http.go` | `MapErrorToHTTPStatus(err)` translation from semantic errors to HTTP status |

## Public API Reference

### Constructors and Helpers

| Function | Returns | Notes |
| --- | --- | --- |
| `ErrMissingUserID()` | `error` | Missing user id in context |
| `ErrUnauthorized(reason)` | `error` | Unauthorized with optional reason |
| `ErrForbidden(reason)` | `error` | Forbidden with optional reason |
| `NewValidationError(field, reason)` | `error` | Validation error with field context |
| `NewAuthenticationError(reason)` | `error` | Authentication failure, mapped as unauthorized |
| `AtLeastOneFieldRequired(fields...)` | `error` | Validation helper for partial update commands |
| `MissingFields(fields...)` | `error` | Validation helper for required field checks |
| `MapErrorToHTTPStatus(err)` | `int` | Main package entrypoint used by HTTP response layer |

### Typed Errors

| Type | Typical status |
| --- | --- |
| `ValidationError` | `400 Bad Request` |
| `UnauthorizedError` | `401 Unauthorized` |
| `ForbiddenError` | `403 Forbidden` |
| `MissingUserIDError` | `401 Unauthorized` |
| `AuthenticationError` | `401 Unauthorized` |

### Sentinel Errors

| Sentinel | Typical status |
| --- | --- |
| `ErrParseUserID` | `400 Bad Request` |
| `ErrNoFieldsToUpdate` | `400 Bad Request` |
| `ErrUsernameExists` | `409 Conflict` |
| `ErrEmailExists` | `409 Conflict` |
| `ErrDomainConflict` | `409 Conflict` |

## HTTP Mapping Summary

| Condition | Status code |
| --- | --- |
| `nil` error | `200 OK` |
| Validation errors | `400` |
| Unauthorized/authentication errors | `401` |
| Forbidden errors | `403` |
| `httperrors.ErrResourceNotFound` | `404` |
| `httperrors.ErrMethodNotAllowed` | `405` |
| Conflict errors | `409` |
| Any unknown error | `500` |

## Usage Example

```go
if err != nil {
    status := sharederrors.MapErrorToHTTPStatus(err)
    httpresponse.WriteError(ctx, w, logger, err)
    _ = status // status can be used for metrics/log decoration if needed
    return
}
```

## Design Notes

- Keep this package transport-semantic only (no domain orchestration).
- Prefer `errors.As` for typed errors and `errors.Is` for sentinels.
- Reuse constants from `internal/shared/constants/commonkeys` when composing field-based validation errors.

## Package Improvements

- Consolidate duplicated semantics for missing user id (`ErrMissingUserID`, `ErrMissingUserIDParam`) into a single source of truth.
- Evaluate whether `MissingFields()` should explicitly handle `len(fields) == 0` to avoid ambiguous validation messages.
- Consider exposing a small test table in package docs to lock the expected map (`error -> status`) behavior.
- Standardize naming between constructor-style functions (`NewValidationError`) and `Err*` helpers for API consistency.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../../README.md)
<!-- doc-nav:end -->
