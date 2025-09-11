# Platform HTTP — Generic (Health / Errors / 404 / 405)

**Folder:** `internal/platform/server/http/generic`
**Subpackages:** `dto/`, `handler/`

## Responsibility

* Provide **platform-wide HTTP controllers** that every context can rely on:

    * `/health` endpoint (service metadata).
    * Central **error handler** (uniform 5xx response + logging).
    * **Recovery** from panics with safe 500 output.
    * Standard **404 Not Found** and **405 Method Not Allowed** responses.
* Keep these cross-cutting concerns **out of domain adapters**, so primary adapters stay thin.

## How it works

* The Generic **Handler** is built with the platform logger and basic app metadata:

    * `New(logger, generalCfg)` → returns `*handler.Handler` (fields: `Logger`, `GeneralConfig`).
* The **server composer** wires it into the router port:

    * `GET /health` → `Handler.HealthCheck`.
    * `SetNotFound(Handler.NotFoundHandler)`.
    * `SetMethodNotAllowed(Handler.MethodNotAllowedHandler)`.
    * `SetErrorHandler(Handler.ErrorHandler)`.
    * Recovery middleware uses `Handler.Recovery` to trap panics and convert to 500.
* All endpoints emit **OpenTelemetry spans** and structured logs (trace names in `0_generic_handler_constants.go`).

## Endpoints / Behaviors

* `GET /health` — Returns service metadata (`name`, `env`, `version`, `timestamp`).
* **404 Not Found** — Standard JSON body, includes `x-request-id` when present.
* **405 Method Not Allowed** — Standard JSON body; do not leak internals.
* **Error handler (500)** — Uniform body + OTel span marked as error; logs with request ID.
* **Recovery (panic → 500)** — Captures stack, tags span as error, returns safe payload.

> All responses follow the shared helpers (e.g., `httpresponse`) to keep the envelope consistent.

## DTOs

**Folder:** `internal/platform/server/http/generic/dto`

* `HealthCheckResponse`

    * `Name string`
    * `Env string`
    * `Version string`
    * `Timestamp time.Time`

Rules of thumb:

* Keep DTOs **transport-only**; no domain types here.
* Add fields **only** if they’re stable across environments and safe to expose.

## Observability

* Tracers defined in constants:
  `aionapi.generic.handler`, `generic.health_check`, `generic.error_handler`, `generic.recovery_handler`.
* Each controller:

    * Starts a span; sets canonical attributes (e.g., request ID).
    * Marks span status `Error` on failures.
    * Logs **metadata** (never sensitive payloads).

## Router Port & Middleware

* Everything is mounted via the **router port** (`internal/platform/server/http/ports.Router`).
* Platform middlewares typically applied in the composer:

    * **Request ID** — guarantees `X-Request-ID` header + context value.
    * **Recovery** — wraps the stack with `Handler.Recovery`.

## Controller conventions

* **No domain logic** here—only platform concerns.
* Prefer returning via shared helpers (keeps JSON envelopes uniform).
* Avoid leaking internal errors; keep messages user-safe, log details separately.
* Never import a concrete router here; depend on the **router port** only.

## Testing hints

* Use `httptest` with the **chi adapter** (or any adapter implementing the port) to exercise real routing:

    * `/health` returns 200 and valid `HealthCheckResponse`.
    * 404/405 return the standardized error format.
    * Error/Recovery paths tag spans as error (when OTel test exporter is configured).
* Assert that `X-Request-ID` is preserved/propagated by the handlers.

## Extending

* Add new **generic** endpoints only when truly cross-cutting (maintenance mode, readiness/liveness, etc.).
* If an endpoint belongs to a **bounded context**, put it in that context’s primary adapter, not here.
* Keep this package **small, deterministic, and framework-agnostic**.
