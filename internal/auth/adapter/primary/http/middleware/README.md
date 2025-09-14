**`internal/auth/adapter/primary/http/middleware/README.md`**

````md
# Auth HTTP Middleware

**Folder:** `internal/auth/adapter/primary/http/middleware`

## Responsibility
* Protect **protected** HTTP routes by validating credentials (e.g., Bearer/JWT) through the **AuthService input port**.
* On success, enrich the `context.Context` with `user_id`, `token`, and optional `claims` for downstream handlers.
* Emit OpenTelemetry spans/attributes and structured logs for both success and failure paths.

## How it works
* Constructed with `New(service, logger)` where:
  - `service` implements `input.AuthService` (token verification only).
  - `logger` is a `logger.ContextLogger`.
* The middleware extracts credentials (typically `Authorization: Bearer <token>`), asks the service to **verify**, and:
  - **Success** → stores auth data in context (e.g., `ctxkeys.UserID`) and calls `next`.
  - **Failure** → short-circuits with **401** and a standardized error body.
* Observability:
  - **Tracer name:** `TracerAuthMiddleware`
  - **Span:** `SpanAuthMiddleware`
  - Error labels: `SpanErrorMissingToken`, `SpanErrorTokenInvalid`
  - Attributes: `AttrAuthMiddlewareStatus`, `AttrAuthMiddlewareUserID`, `AttrAuthMiddlewareError`

> No persistence/ORM work here; the middleware talks only to the **input port**.

## Usage
Register protected subtrees using the router **port** (`ports.Router`) and `GroupWith`:

```go
r.Group("/v1/users", func(gr ports.Router) {
  gr.GroupWith(middleware.Auth(service, logger), func(pr ports.Router) {
    pr.GET("/{user_id}", http.HandlerFunc(h.GetByID))
    pr.PUT("/",           http.HandlerFunc(h.UpdateUser))
    pr.DELETE("/",        http.HandlerFunc(h.SoftDeleteUser))
  })
})
````

Keep public endpoints outside the `GroupWith(...)` block.

## Context contract (downstream)

* Read authenticated data via `ctxkeys` (e.g., `ctxkeys.UserID`) in handlers.
* Do **not** re-parse tokens in handlers; rely on the middleware to populate context.
* Never log or echo the raw token; log **metadata only**.

## Error handling conventions

* Missing/empty token → **401 Unauthorized** with a canonical error envelope.
* Invalid/expired token → **401 Unauthorized** (do not leak verification details).
* The middleware records errors on the span and logs metadata (no secrets).

## Design notes

* Transport-agnostic: depends on `ports.Router` and runs as an HTTP middleware.
* Security: prefer **Authorization: Bearer**. If cookies are added later, keep the same verification path in the service.
* Use shared helpers (`httpresponse`, `sharederrors`) for consistent error bodies.

## Testing hints

* Stub `AuthService.Validate/Verify` to return:

    * **Success** → assert `next` was called and `ctxkeys.UserID` is present.
    * **Failure** → assert **401** and `next` was **not** called.
* Cover edge cases:

    * Missing header / malformed `Bearer` token
    * Context canceled before verification
    * Service error vs. explicit “invalid token”
* Tracing/logging: relax logger expectations in gomock (`AnyTimes()`) and, in integration-style tests, assert span attributes when an OTEL test exporter is configured.

## Tracing (OTel)

* Use `TracerAuthMiddleware`; start a span per request (`SpanAuthMiddleware`).
* Set attributes:

    * `AttrAuthMiddlewareStatus`: `"authenticated"` on success, or an error code on failure
    * `AttrAuthMiddlewareUserID`: present on success
    * `AttrAuthMiddlewareError`: brief reason tag on failures
* On errors: `span.RecordError(err)` and set an error status; avoid secrets in events.


