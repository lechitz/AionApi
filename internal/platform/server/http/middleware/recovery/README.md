# HTTP Middleware — Recovery (Platform)

**Folder:** `internal/platform/server/http/middleware/recovery`

## Responsibility

* Catch **panics** anywhere in the HTTP pipeline and convert them into a safe **500 Internal Server Error**.
* Guarantee a **request ID** on the response (generate one if missing) to aid correlation.
* Delegate response formatting + logging to the platform **Generic Recovery Handler**.
* Preserve **observability**: mark spans as error and attach useful attributes.

## How it works

* Constructor: `New(recoveryHandler *generic.Handler) func(http.Handler) http.Handler`

    * Wraps `next` with a `defer`/`recover()` guard.
    * If a panic occurs:

        * Ensures an `X-Request-ID` (creates a UUID when absent).
        * Calls the platform **Generic Recovery** controller to:

            * log with structure, include request ID
            * capture stack (internally)
            * emit an OTel span with `status=Error`
            * return a standardized 500 JSON body.
* Depends only on the **router port** middleware signature (`func(http.Handler) http.Handler`).

## Usage

Apply globally in the HTTP composer so it wraps **all** routes and middlewares:

```go
// composer.go
r.Use(recovery.New(genericHandler))  // <- outermost guard
```

> Order matters: keep Recovery **first** in the chain so it can catch panics from every subsequent middleware/handler.
> If you also use `requestid` middleware, it can come after Recovery—Recovery will still generate a fallback request ID if needed.

## Observability

* The Recovery handler:

    * starts/annotates an OTel span for the failing request
    * sets canonical attributes (e.g., request ID)
    * marks span `Error` and records the panic message/stack (without leaking it to clients).

## Conventions

* Keep client messages **generic** (no internal details). Detailed error context goes to logs/spans.
* Never perform domain logic here—this is a **platform** concern only.
* Do not import concrete routers; wire via the **router port**.

## Testing hints

* Build a tiny handler that `panic`s and mount it behind the middleware:

    * Assert response `500`, JSON envelope shape, and presence of `X-Request-ID`.
    * Verify no server crash and that logs/spans mark the request as error.
* When combined with `requestid` middleware, assert the **same** request ID flows through.
