# Shared Errors

**Folder:** `internal/shared/sharederrors`

## Responsibility
* Provide **reusable semantic error types** (e.g., NotFound, Unauthorized, Conflict) used across use cases and adapters.
* Enable **uniform HTTP mapping** via `internal/shared/httpresponse` so handlers don’t reinvent status codes.

## How it works
* Exposes constructors/helpers for common categories:
  - **Validation** (bad input, missing fields)
  - **Unauthorized / Forbidden**
  - **NotFound**
  - **Conflict / AlreadyExists / Uniqueness**
  - **TooManyRequests**
  - **Internal / Unexpected**
* Each error carries a **stable code** and a human-safe message. Transport layers map codes to HTTP status.

## Typical mapping (by `httpresponse`)
* `validation` → **400 Bad Request**
* `unauthorized` → **401 Unauthorized**
* `forbidden` → **403 Forbidden**
* `not_found` → **404 Not Found**
* `conflict` → **409 Conflict**
* `rate_limited` → **429 Too Many Requests**
* `internal` → **500 Internal Server Error**

## Usage
```go
// In a use case:
if username == "" {
  return sharederrors.Validation("username is required")
}

// In a repository when a row is missing:
return sharederrors.NotFound("user not found")
````

```go
// In an HTTP handler:
if err != nil {
  // Centralized mapping + standardized envelope
  httpresponse.WriteError(ctx, w, logger, err)
  return
}
```

## Design notes

* Keep messages **non-sensitive**; log metadata only.
* Prefer these errors over ad-hoc strings to keep telemetry and responses consistent.

## Testing hints

* Assert **error codes**, not string messages.
* Table-test the HTTP mapping: each error → expected status and envelope.

```

Say the word and I’ll send the next one.
::contentReference[oaicite:0]{index=0}
```
