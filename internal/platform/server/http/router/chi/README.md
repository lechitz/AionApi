# HTTP Router Adapter — chi

**Folder:** `internal/platform/server/http/router/chi`

## Responsibility

* Implement the platform routing contract **`internal/platform/server/http/ports.Router`** using **`github.com/go-chi/chi/v5`**.
* Keep application contexts (Auth/User/Category/…) **decoupled** from the concrete router.
* Provide a small, consistent API for grouping, mounting, defaults, and middleware.

---

## How it works

* The adapter wraps a `chi.Router` in a private type that **implements `ports.Router`**.

* It exposes the same surface expected by contexts and the HTTP composer:

    * **Middlewares:** `Use(mw ...Middleware)`
    * **Composition:** `Group`, `GroupWith`, `Mount`
    * **Handlers:** `Handle`, `GET`, `POST`, `PUT`, `DELETE`
    * **Defaults:** `SetNotFound`, `SetMethodNotAllowed`, `SetError`
    * **Entry point:** `ServeHTTP`

* `SetError` stores a function used by the **platform error/recovery path** (the framework will not call it automatically). The platform recovery middleware invokes it to render uniform 500s.

* Constructor: `New() ports.Router`

---

## Usage (HTTP composer excerpt)

```go
import (
  chiadapter "github.com/lechitz/AionApi/internal/platform/server/http/router/chi"
  "github.com/lechitz/AionApi/internal/platform/server/http/middleware/requestid"
  "github.com/lechitz/AionApi/internal/platform/server/http/middleware/recovery"
  generic "github.com/lechitz/AionApi/internal/platform/server/http/generic/handler"
  userhandler "github.com/lechitz/AionApi/internal/user/adapter/primary/http/handler"
)

func ComposeHTTP(...) ports.Router {
  r := chiadapter.New() // concrete adapter that implements ports.Router

  // Global middlewares (order matters)
  r.Use(recovery.New(genericHandler)) // outermost guard
  r.Use(requestid.New())              // ensure X-Request-ID

  // Platform defaults
  r.SetNotFound(http.HandlerFunc(genericHandler.NotFoundHandler))
  r.SetMethodNotAllowed(http.HandlerFunc(genericHandler.MethodNotAllowedHandler))
  r.SetError(genericHandler.ErrorHandler)

  // Context routes
  userhandler.RegisterHTTP(r, userH)
  userhandler.RegisterHTTPProtected(r, userH, authMiddleware)

  // GraphQL (example)
  r.Mount(cfg.ServerGraphql.Path, graphqlHTTPHandler)

  return r
}
```

---

## Grouping & protected subtrees

```go
// Public subtree
r.Group("/v1/users", func(ur ports.Router) {
  ur.POST("/create", http.HandlerFunc(h.Create))
})

// Protected subtree (middleware applies only here)
r.Group("/v1/users", func(ur ports.Router) {
  ur.GroupWith(authMiddleware, func(pr ports.Router) {
    pr.GET("/{user_id}", http.HandlerFunc(h.GetByID))
    pr.PUT("/",           http.HandlerFunc(h.UpdateUser))
    pr.PUT("/password",   http.HandlerFunc(h.UpdatePassword))
    pr.DELETE("/",        http.HandlerFunc(h.SoftDelete))
  })
})
```

---

## Mounting composite handlers

Prefer `Mount` for fully built handlers (e.g., GraphQL):

```go
r.Mount(cfg.ServerGraphql.Path, graphqlHTTPHandler)
```

---

## Observability

* The adapter is intentionally **thin**; tracing/logging live in:

    * Platform middlewares (`requestid`, `recovery`)
    * Platform generic handlers (health, 404/405, error)
    * Context handlers (User/Auth/Category), which start OTel spans and log with structured fields
* Always apply `requestid` early to guarantee correlation across logs/traces/responses.

---

## Swapping adapters

To change the routing engine, provide another implementation under `internal/platform/server/http/router/<impl>` that satisfies `ports.Router`, and switch the import in the composer:

```go
r := otheradapter.New() // still returns ports.Router
```

No changes are needed in context registrars.

---

## Testing hints

* Use `httptest` with this adapter to exercise real routing while keeping contexts bound to `ports.Router`.
* Assert:

    * Defaults: 404/405 handlers are used
    * Middleware order effects (recovery catches panics; request ID is present)
    * Protected subtrees do not allow access without the auth middleware
* For unit tests of registrars, inject a **fake `ports.Router`** and assert calls.

---

## Notes / Gotchas

* **Middleware order matters** — place `recovery` first so it can catch panics from everything else.
* `GroupWith` applies the given middleware **only** to that subtree.
* `SetError` does not hook into chi automatically; it’s invoked by the platform **recovery/error** flow.
* Context packages must **not** import `chi`; depend only on `internal/platform/server/http/ports`.
