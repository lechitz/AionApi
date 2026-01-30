# HTTP Utils — Shared Errors

**Folder:** `internal/platform/server/http/utils/sharederrors`

## Purpose and Main Capabilities

- Provide semantic error helpers used across handlers and adapters.
- Enable uniform HTTP status mapping via `httpresponse`.
- Keep error codes stable and client-safe.

## How it works

- Exposes constructors for common categories:
  - Validation
  - Unauthorized / Forbidden
  - NotFound
  - Conflict / AlreadyExists
  - TooManyRequests
  - Internal

## Typical mapping (via `httpresponse`)

| Error code | HTTP status |
| --- | --- |
| `validation` | 400 |
| `unauthorized` | 401 |
| `forbidden` | 403 |
| `not_found` | 404 |
| `conflict` | 409 |
| `rate_limited` | 429 |
| `internal` | 500 |

## Usage

```go
if username == "" {
  return sharederrors.Validation("username is required")
}
```

```go
if err != nil {
  httpresponse.WriteError(ctx, w, logger, err)
  return
}
```

## Conventions

- Keep messages non-sensitive.
- Prefer semantic errors over ad-hoc strings.

## What Should NOT Live Here

- Domain logic or transport DTOs.
```
