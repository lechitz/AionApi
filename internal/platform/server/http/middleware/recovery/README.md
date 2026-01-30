# HTTP Middleware — Recovery (Platform)

**Folder:** `internal/platform/server/http/middleware/recovery`

## Purpose and Main Capabilities

- Catch panics in the HTTP pipeline and return a safe 500 response.
- Delegate formatting/logging to the generic recovery handler.
- Preserve observability by marking spans as error.

## How it works

- `New(recoveryHandler *handler.Handler)` wraps `next` with `defer`/`recover()`.
- On panic, it generates an error ID (UUID) and calls `recoveryHandler.RecoveryHandler`.

## Usage

Apply globally as the **outermost** middleware:

```go
r.Use(recovery.New(genericHandler))
```

## Observability

- Recovery handler records error details on spans and logs.
- Client receives a sanitized 500 response only.

## What Should NOT Live Here

- Domain logic or transport mapping.
