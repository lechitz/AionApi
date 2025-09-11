# Generic HTTP Handlers (Platform)

**Path:** `internal/platform/server/http/generic/handler`

## Purpose

Platform-wide controllers for common HTTP concerns: health check, standardized 404/405 responses, and centralized error/panic handling. Keeps these cross-cutting behaviors out of domain adapters.

## Responsibilities

* `/health` endpoint returning service metadata and status.
* Standard 404 (Not Found) and 405 (Method Not Allowed) responses.
* Central **ErrorHandler** used by the platform router as a default error sink.
* **RecoveryHandler** to turn panics into a safe 500 response (wired by recovery middleware).

## Why it matters

* Eliminates duplication across contexts (Auth/User/…).
* Ensures **consistent logs, traces, and payloads** for common responses.
* Provides a single place to evolve platform behavior (e.g., headers, shapes, telemetry).

## How it works

* **Constructor:** `New(logger, generalCfg)` injects the contextual logger and app metadata.
* **Wire in composer:**

    * Register defaults: `SetNotFound`, `SetMethodNotAllowed`, `SetError`.
    * Add recovery middleware with `RecoveryHandler` (gets an `errorID` for correlation).
* **Responses:** Written via `internal/shared/httpresponse` with stable shapes/keys.
* **DTOs:** Health payload uses `generic/dto.HealthCheckResponse`.

## Observability

* OTel spans for health, error, and recovery paths.
* Structured logs with canonical keys (request id, path, method, IP, user agent, error id).
* No domain attributes or PII are logged.

## Integration (excerpt)

```go
gh := generic.New(log, cfg.General)

// middlewares (order matters)
r.Use(requestid.New(), recovery.New(gh))

// defaults
r.SetNotFound(http.HandlerFunc(gh.NotFoundHandler))
r.SetMethodNotAllowed(http.HandlerFunc(gh.MethodNotAllowedHandler))
r.SetError(gh.ErrorHandler)

// route
api.GET("/health", http.HandlerFunc(gh.HealthCheck))
```

## Non-goals

* No domain/business logic.
* No transport coupling beyond `net/http`.
* No DB/cache access.

## Notes / Tips

* Keep messages and keys stable for dashboards and alerts.
* Extend here for additional platform endpoints (e.g., `/ready`, security headers) if needed.
* Unit test handlers directly with `httptest`; avoid spinning up a full server.
