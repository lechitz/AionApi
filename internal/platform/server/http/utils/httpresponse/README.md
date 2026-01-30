# HTTP Utils — Response Helpers

**Folder:** `internal/platform/server/http/utils/httpresponse`

## Purpose and Main Capabilities

- Centralize JSON response formatting.
- Map semantic errors to HTTP status codes.
- Provide a consistent error envelope across adapters.

## How it works

- Helpers write JSON responses and set `Content-Type`.
- `WriteError` (or equivalent) maps `sharederrors` to status codes.
- Structured logging uses `logger.ContextLogger` with `commonkeys`.

## Error mapping (typical)

- validation -> 400
- unauthorized -> 401
- forbidden -> 403
- not_found -> 404
- conflict -> 409
- rate_limited -> 429
- internal -> 500

## Usage (pattern)

```go
if err != nil {
  httpresponse.WriteError(ctx, w, logger, err)
  return
}
```

## Conventions

- Always return JSON; never leak sensitive data.
- Prefer semantic errors from `sharederrors`.

## What Should NOT Live Here

- Domain logic or handler orchestration.
