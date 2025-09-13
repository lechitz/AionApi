# User HTTP (Primary Adapter)

**Folder:** `internal/user/adapter/primary/http/handler`

## Responsibility

* Expose the **User** context over HTTP (controllers + route registration).
* Keep handlers **thin**: decode/validate inputs, call **input ports**, map responses, and standardize errors.
* Emit **observability** signals (OTel spans + structured logs), but **no business logic** here.

## How it works

* **Dependencies** are injected with `New(userService, cfg, logger)` (see `Handler` fields: `UserService`, `Logger`, `Config`).
* Route groups are registered through the **router port** (`platform/server/http/ports.Router`) to avoid coupling to a concrete engine.
* The Platform composer mounts everything under the global API prefix (e.g., `/aion-api`).
* Shared helpers:

    * **`handlerhelpers`** + **`httpresponse`** for consistent error/JSON responses.
    * **`httputils`** for auth-cookie handling when relevant (e.g., password updates / signout).
    * **`sharederrors`** to convert domain errors into HTTP status codes/messages.

## Routes

### Public (`RegisterHTTP`)

* `POST /v1/users/create` — create a new user (only if your product allows open sign-up).

### Protected (`RegisterHTTPProtected` + `AuthMiddleware`)

* `GET    /v1/users/all` — list users (**RBAC recommended: admin-only**).
* `GET    /v1/users/{user_id}` — fetch a user by ID (**self or admin**).
* `PUT    /v1/users/` — update the authenticated user (partial updates supported).
* `PUT    /v1/users/password` — update the authenticated user’s password (may refresh cookie).
* `DELETE /v1/users/` — soft delete the authenticated user (clears cookie).

> Effective paths include the global API prefix, e.g. `/aion-api/v1/users/...`.

## Controller conventions

* Start an **OTel span** per handler; use constants from `0_user_handler_constants.go` and add canonical attributes (e.g., `user_id`, `http.status_code`).
* **Never** touch persistence/ORM here; call **input ports** on `UserService`.
* Validate at the boundary and map DTOs in/out; **do not** expose domain entities directly.
* Prefer **semantic domain errors** (via `sharederrors`) over ad-hoc strings.
* Log **metadata** only—avoid sensitive values (passwords, tokens, cookie contents).

## Files overview

* `0_user_handler_constants.go` — tracer/span names and shared errors.
* `0_user_handler_impl.go` — `Handler` struct (DI surface).
* `0_user_handler_helpers.go` — small helpers (param parsing, common response flows).
* `create_user_handler.go` — controller for `POST /create`.
* `getall_user_handler.go` — controller for `GET /all`.
* `get_by_id_user_handler.go` — controller for `GET /{user_id}`.
* `update_user_handler.go` — controller for `PUT /` (partial update).
* `update_password_user_handler.go` — controller for `PUT /password`.
* `softdelete_user_handler.go` — controller for `DELETE /`.
* `register.go` / `register_protected.go` — route registration into the `ports.Router`.

## Design notes

* Use the **router port** only; concrete routers live under `internal/platform/server/http/router/*`.
* Apply auth only to protected subtrees using `GroupWith(middleware.Auth, ...)`.
* Keep handlers **small and deterministic**—easy to unit test with stubbed services.

## Testing hints

* Unit test each handler with a **fake `UserService`**:

    * assert **status codes** and **response bodies**
    * validate **input errors** and missing/invalid params
    * ensure **partial updates** map only provided fields
    * verify **error propagation** and standardized error mapping
* For cookies, assert `Set-Cookie` headers (HttpOnly, Secure, SameSite, Max-Age/Expires) when applicable.
