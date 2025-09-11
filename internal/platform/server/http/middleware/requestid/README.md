# HTTP Middleware — Request ID (Platform)

**Folder:** `internal/platform/server/http/middleware/requestid`

## Responsibility

* Ensure **every request** carries a stable **`X-Request-ID`**.
* Accept client-provided IDs (if present & non-empty) or **generate a UUIDv4**.
* Propagate the ID to:

    * the **request context** (for logs/traces/handlers), and
    * the **response header** (so clients can correlate).
* Stay framework-agnostic via the **router port** middleware signature.

## How it works

* Constructor: `New() func(http.Handler) http.Handler`

    * Reads `X-Request-ID` (header name from `commonkeys.XRequestID`).
    * If missing/blank → `uuid.NewString()`.
    * Stores in context under `ctxkeys.RequestID` and sets the header on both request (for downstream middlewares/handlers) and response.
    * Calls `next.ServeHTTP` with the updated context.

> No domain logic here—pure platform concern for correlation and observability.

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

* Prefer **UUIDv4** when generating new IDs.
* Treat client-provided IDs as opaque values—**do not parse** or validate beyond non-empty.
* Keep the middleware **stateless** and **idempotent**.

## Testing hints

* With a pre-set `X-Request-ID`, assert the same value is present in:

    * request context (`ctxkeys.RequestID`),
    * response header.
* Without a header, assert a **non-empty UUID-like** value is generated and propagated consistently.
* Combine with the **Recovery** middleware and ensure the same ID appears in 500 responses/logs after a panic.
