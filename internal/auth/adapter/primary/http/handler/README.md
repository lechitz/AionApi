# Auth HTTP (Primary Adapter)

**Folder:** `internal/auth/adapter/primary/http/handler`

## Responsibility

* Expose the **Auth** context over HTTP (controllers + route registration).
* Keep controllers **thin**: decode/validate input, call input ports, map responses, standardize errors.
* Manage session lifecycle concerns at the edge (cookies/headers) without leaking into the domain.

## How it works

* Dependencies are injected with `New(service, cfg, logger)`.
* Route groups are registered via the **router port** (`ports.Router`) to avoid coupling to a concrete engine.
* Protected subtrees use the auth domain middleware (`middleware.Auth`) with `GroupWith`.
* Handlers start an OTel span using the tracer name in `TracerAuthHandler` and span names like `auth.login`, `auth.logout`.

> Platform composition may mount everything under a global prefix (e.g., `/aion-api`).

## Routes

### Public (`RegisterHTTP`)

* `POST /v1/auth/login` — authenticate user and create a session.

### Protected (`RegisterHTTP` + `middleware.Auth`)

* `POST /v1/auth/logout` — terminate the authenticated user session.

> Effective paths include any global API prefix, e.g., `/aion-api/v1/auth/...`.

## Controller conventions

* Use shared helpers (`httpresponse`, `sharederrors`, `httputils`) to normalize responses and headers.
* Add canonical OTel attributes (e.g., `operation`, `user_id` when available) via `tracingkeys`.
* **Never** log secrets: passwords, raw tokens, or cookie values. Log **metadata only**.
* Do not touch persistence/ORM here; call the **input ports** on the Auth service.
* Read authenticated context (e.g., user id/claims) via `ctxkeys` in protected handlers.
* Centralize cookie/header settings (e.g., session token) using helper utilities—keep policy in the service.

---

## DTOs

**Folder:** `internal/auth/adapter/primary/http/dto`

### Purpose

* Define **request/response shapes** for HTTP only.
* Perform **boundary validation** (required fields, formats).
* Map to **input-port commands** (keeps the domain clean/testable).

### Typical types

* `LoginRequest` → maps to `input.AuthService.Login(...)`.
* `LoginResponse` — session material (e.g., sets cookie and/or returns a short response model).
* `LogoutResponse` — standardized success envelope.

### Rules of thumb

* Validate at the boundary (username/email format, password presence).
* **Never** echo back sensitive fields.
* Mapping is one-way: DTO → command (input port), and domain → DTO for outputs.

---

## Design notes

* Stick to the router **port** (`ports.Router`); avoid importing concrete routers.
* Apply `middleware.Auth` only where needed.
* Use `httpresponse` helpers to ensure consistent status codes and error bodies.
* Keep handlers deterministic and small—easy to unit test with a stubbed service.

## Testing hints

* Unit test each handler with a fake `AuthService`, asserting:

    * Correct status codes and response envelopes.
    * Validation failures (missing/invalid credentials).
    * Cookie/header behavior on login/logout.
    * Propagation of domain errors (e.g., invalid credentials → 401).
    * Tracing/logging hooks are called (optional, usually relaxed with gomock `AnyTimes()`).

## Tracing (OTel)

* Tracer: `TracerAuthHandler`.
* Spans: `auth.login`, `auth.logout`.
* Attributes to consider: `operation`, `user_id` (when authenticated), `http_status`, and **status** tags for success/failure.
* On errors: `span.RecordError(err)` and set an error status; keep logs free of sensitive data.
