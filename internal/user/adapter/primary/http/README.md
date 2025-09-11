# User HTTP (Primary Adapter)

**Folder:** `internal/user/adapter/primary/http`

## Responsibility

* Expose the **User** context over HTTP (controllers + route registration).
* Keep controllers **thin**: decode/validate, call input ports, map responses, and standardize errors.

## How it works

* Dependencies are injected with `New(userService, cfg, logger)`.
* Route groups are registered via the **router port** (`ports.Router`) to avoid coupling to a concrete engine.
* The Platform composer mounts everything under `cfg.ServerHTTP.Context` (e.g., `/aion-api`).

## Routes

### Public (`RegisterHTTP`)

* `POST /v1/users/create` — create a new user (use only if your product allows open sign-up).

### Protected (`RegisterHTTPProtected` + `AuthMiddleware`)

* `GET    /v1/users/all` — list users (**RBAC recommended: admin-only**).
* `GET    /v1/users/{user_id}` — fetch a user by ID (**self or admin**).
* `PUT    /v1/users/` — update the authenticated user.
* `PUT    /v1/users/password` — update the authenticated user’s password (refresh cookie).
* `DELETE /v1/users/` — soft delete the authenticated user (clears cookie).

> Effective paths include the global API prefix, e.g. `/aion-api/v1/users/...`.

## Controller conventions

* Use `handlerhelpers`, `httpresponse`, and `sharederrors` to standardize responses.
* Start an OTel span per handler; add canonical attributes (user\_id, operation, http\_status).
* Never access persistence/ORM here; call the **input ports** on `UserService`.
* Prefer **semantic errors** from the domain instead of ad-hoc strings.

---

## DTOs

**Folder:** `internal/user/adapter/primary/http/dto`

### Purpose

* Define **request/response shapes** for HTTP only.
* Perform **input validation** close to the transport layer.
* Map to **input-port commands** (keeps domain clean and testable).

### Key types

* `CreateUserRequest` (+ `ValidateUser`) → maps to `input.CreateUserCommand`.
* `CreateUserResponse` — HTTP output after creation.
* `UpdateUserRequest` → maps to `input.UpdateUserCommand` (optional fields; only provided fields update).
* `UpdateUserResponse` — updated fields + `updated_at` + `user_id`.
* `UpdatePasswordUserRequest` — current & new password.
* `GetUserResponse` — read model for `GET` endpoints.

### Rules of thumb

* **Validate** at the boundary (e.g., required fields, email format, password length).
* **Never** expose domain entities directly; always map to DTOs.
* Keep mapping **one-way**: DTO → command (input port), and domain → DTO for responses.

---

## Design notes

* Stick to the router **port** (`ports.Router`); avoid importing concrete routers in new code.
* Apply auth only where needed with `GroupWith(middleware.Auth, ...)`.
* Log **metadata**, not sensitive payloads (passwords, tokens, etc.).
* Handlers should be small and deterministic—easy to unit test with stubbed services.

## Testing hints

* Unit test each handler with a fake `UserService`, asserting:

    * correct status codes/messages
    * validation failures
    * mapping of partial updates
    * error propagation from the service.
