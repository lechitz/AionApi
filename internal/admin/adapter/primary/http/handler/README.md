# Admin HTTP (Primary Adapter)

**Folder:** `internal/admin/adapter/primary/http/handler`

## Responsibility

* Expose the **Admin** context over HTTP (handlers + route registration hooks).
* Keep handlers **thin**: decode/validate, call input ports, map responses, and standardize errors.
* Enforce **auth & RBAC** (admin-only) at the boundary.

## How it works

* Dependencies are injected via a constructor like `New(adminService, cfg, logger)` in the adapter layer (handler files stay transport-only).
* Route groups are registered through the **router port** (`ports.Router`) to avoid coupling to a specific engine.
* The Platform composer mounts everything under `cfg.ServerHTTP.Context` (e.g., `/aion-api`).
* Each handler starts an **OTel span** and sets canonical attributes (operation, user\_id, http\_status).

## Routes (suggested layout)

### Public (rare)

* Typically **none**. Admin endpoints should be protected.

### Protected (`RegisterHTTPProtected` + `AuthMiddleware` + **RBAC: admin**)

* `GET    /v1/admin/health` — service health / dependencies readiness.
* `GET    /v1/admin/users` — list users with admin filters (pagination recommended).
* `PUT    /v1/admin/users/{user_id}` — admin update of user profile/roles.
* `DELETE /v1/admin/users/{user_id}` — admin soft-delete/restore.
* `GET    /v1/admin/stats` — platform metrics (aggregations, counters).
* `POST   /v1/admin/feature-flags` — toggle/assign flags (if applicable).

> Effective paths include the global API prefix, e.g. `/aion-api/v1/admin/...`.

## Handler conventions

* Use `handlerhelpers`, `httpresponse`, and `sharederrors` to **normalize** responses and errors.
* Do **not** reach infra here (no ORM/SQL/Redis). Call the **input ports** on the Admin service.
* Log **metadata only** (no PII/secrets). Prefer structured `...wCtx` logging with request-scoped context.
* Map **domain errors → HTTP** consistently (e.g., `ErrNotFound → 404`, `ErrUnauthorized → 401/403`, `ValidationError → 422`).

---

## DTOs

**Folder (suggested):** `internal/admin/adapter/primary/http/dto`

### Purpose

* Define **request/response** shapes for HTTP only.
* Perform **input validation** close to the transport layer.
* Map to **input-port commands** (keep domain clean/testable).

### Key types (examples)

* `UpdateUserAsAdminRequest` → maps to `admininput.UpdateUserCommand`.
* `AdminListUsersRequest` (query params for filters/pagination).
* `AdminUserResponse`, `AdminStatsResponse`, `AdminHealthResponse`.

### Rules of thumb

* **Validate** at the boundary (required fields, enums, pagination bounds).
* Never expose domain entities directly; always map to DTOs.
* Keep mapping **one-way**: DTO → command (input port), and domain → DTO for responses.

---

## Design notes

* Stick to the router **port** (`ports.Router`); avoid importing concrete routers in new code.
* Group routes with `GroupWith(middleware.Auth, middleware.RequireRole(Admin), ...)`.
* Be conservative exposing data—admin endpoints often carry elevated privileges.
* Prefer **idempotent** updates and **explicit** side-effect endpoints.

## Testing hints

* Unit test each handler with a **fake AdminService**, asserting:

    * correct status codes/messages
    * validation failures (bad input, RBAC)
    * mapping of partial updates and pagination
    * error propagation from the service (including domain → HTTP mapping)
* Keep tests transport-focused: build request/response, stub service calls, assert JSON body & headers.

---

## Status

The folder is scaffolded. Implement concrete handlers (e.g., `RegisterHTTP`, `RegisterHTTPProtected`, and per-endpoint functions) following the conventions above and mirroring the Auth/User adapters’ style.
