# HTTP Middleware — Request ID (Platform)

**Folder:** `internal/platform/server/http/middleware/requestid`

## Purpose and Main Capabilities

- Ensure every request carries a stable `X-Request-ID`.
- Validate client-provided IDs and generate a UUID when invalid.
- Propagate the ID to context and response headers.

## How it works

- `New()` reads `X-Request-ID` (`commonkeys.XRequestID`).
- If missing, too long, or not a UUID, it generates a new UUID.
- Injects the final value into `ctxkeys.RequestID` and sets the response header.

## Usage

Register early in your HTTP composer:

```go
// composer.go
r.Use(recovery.New(genericHandler)) // catch panics first
r.Use(requestid.New())              // ensure every request has an ID
```

### Reading the Request ID

```go
func (h *Handler) SomeEndpoint(w http.ResponseWriter, r *http.Request) {
    rid, _ := r.Context().Value(ctxkeys.RequestID).(string)
    h.Logger.Infow("processing request", "request_id", rid)
    // ...
}
```

Clients will also see the same value echoed back in the `X-Request-ID` response header.

## Observability

* The ID becomes a **correlation key** across:

    * structured logs (`request_id` field),
    * tracing (added as an attribute if your handlers/middlewares record it),
    * client–server exchanges (header present on both sides).

## Conventions

- Prefer UUIDv4 when generating new IDs.
- Keep the middleware stateless and idempotent.

## Testing hints

* With a pre-set `X-Request-ID`, assert the same value is present in:

    * request context (`ctxkeys.RequestID`),
    * response header.
* Without a header, assert a **non-empty UUID-like** value is generated and propagated consistently.
* Combine with the **Recovery** middleware and ensure the same ID appears in 500 responses/logs after a panic.
