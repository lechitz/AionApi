# Platform Server

**Folder:** `internal/platform/server`

**Subtrees:** `http/`, `graph/`

## Responsibility

* Provide the platform **transport layer** for HTTP and GraphQL.
* Own the **router port** (`internal/platform/server/http/ports.Router`) and concrete router adapters (e.g., **chi**).
* Offer **cross-cutting handlers** (health check, 404/405, error/recovery) and **middleware** (request ID, panic recovery).
* **Compose** primary adapters from each bounded context (Auth, Category, etc.) into one runnable server.

---

## Layout

```
internal/platform/server
├─ http/
│  ├─ ports/          # framework-agnostic Router contract
│  ├─ router/chi/     # chi implementation of the Router port
│  ├─ generic/        # platform controllers (health, 404/405, error, recovery)
│  ├─ middleware/     # request ID, recovery middlewares
│  ├─ composer.go     # mounts routes from contexts + platform endpoints
│  └─ server.go       # HTTP server startup/shutdown
└─ graph/
   └─ schema/_modules # assembled GraphQL schema (via `make graphql`)
```

---

## How it works

* **Port + Adapter**

    * Contexts (Auth/Category/…) depend only on `http/ports.Router`.
    * We can swap the underlying router by providing a different adapter under `http/router/<impl>`.

* **Composition**

    * `http/composer.go` wires platform endpoints and each context’s **primary HTTP adapter**:

        * mounts **/health** and default **404/405** handlers,
        * calls each context’s `RegisterHTTP(r ports.Router, h *Handler)` (or equivalent),
        * applies platform middlewares (request ID, recovery),
        * mounts everything under `cfg.ServerHTTP.Context` (e.g., `/aion-api`) for a stable API prefix.

* **Server lifecycle**

    * `http/server.go` builds a standard `*http.Server` using timeouts/host/port from `internal/platform/config`.
    * Exposes start/shutdown helpers respecting graceful shutdown timeouts.

* **Observability**

    * Generic handlers and router adapter include **OTel spans** and structured logs through the platform logger.
    * Recovery middleware captures panics, tags spans as errors, and returns a consistent 500 response.
    * Request-ID middleware guarantees a `X-Request-ID` header and context value for traceability.

---

## Router Port (essentials)

`internal/platform/server/http/ports.Router` exposes a minimal, consistent API:

* `Use(mw ...Middleware)` — apply platform middlewares.
* `Group(prefix string, fn func(ports.Router))` — mount subtrees.
* `GroupWith(mw Middleware, fn func(ports.Router))` — mount protected subtrees (e.g., Auth).
* `GET/POST/PUT/PATCH/DELETE(path string, h http.HandlerFunc)` — register handlers.
* Setters for **NotFound/MethodNotAllowed** and a central **ErrorHandler**.

> Context adapters must **only** import this port, never a concrete router package.

---

## GraphQL

* The **schema is modular** per context (e.g., `internal/category/adapter/primary/graphql/schema/*.graphqls`).
* `make graphql`:

    * copies all context schemas into `internal/platform/server/graph/schema/_modules`,
    * runs **gqlgen** to generate GraphQL types/resolvers near the platform graph server,
    * runs `go mod tidy` to keep modules clean.
* The GraphQL HTTP endpoint is composed by the platform alongside REST routes (see the project’s GraphQL server adapter).

---

## Conventions

* **Primary adapters stay thin**: validate/decode, call input ports, map responses, standardize errors.
* Always register routes via the **router port**; never import `chi` (or any router) directly from contexts.
* Apply **domain middleware** (e.g., Auth) with `GroupWith` only where required.
* Log **metadata** (request ID, user ID, operation) — avoid sensitive payloads.

---

## Testing hints

* Use the **chi router adapter** with `httptest` to exercise real HTTP behavior while still going through the port.
* For unit tests of handlers, inject fakes/mocks for input ports and assert:

    * correct status codes and response envelopes,
    * validation failures and error mapping,
    * RBAC/auth middleware behavior when mounted behind `GroupWith`.
* Platform endpoints (health, 404/405, recovery) can be validated with small black-box tests.

---

## Extending

* **New context over HTTP**: create the primary adapter (`RegisterHTTP`, handlers, DTOs), then mount it in `http/composer.go`.
* **New middleware**: implement `ports.Middleware` and add it via `Use` in the composer.
* **Swap router**: implement `ports.Router` under `http/router/<impl>` and switch the binding in the platform bootstrap.
* **Add GraphQL fields**: drop `.graphqls` files in your context module; run `make graphql`; wire resolvers in your context’s GraphQL handler and let the composer mount the server.

> Keep the server package **framework-agnostic, observable, and composable**. The goal is to isolate infrastructure choices while giving every context a uniform, testable entry point.
