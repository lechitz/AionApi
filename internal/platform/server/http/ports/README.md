# Router Port (Platform → HTTP)

**Path:** `internal/platform/server/http/ports`

## Purpose

Defines the **framework-agnostic** routing contract (`Router`) used by all primary HTTP adapters (contexts like `auth`, `user`, etc.). This lets you swap the underlying router (chi/gin/echo) without touching domain code.

## Why it matters

* **Decoupling:** Contexts depend only on `ports.Router`, never on a concrete router package.
* **Replaceability:** Swap routing engines by adding a new adapter under `router/<impl>` that implements this port.
* **Consistency:** All contexts register routes the same way (`RegisterHTTP(r ports.Router, h *Handler)`).

## Contract Overview

```go
type Middleware func(http.Handler) http.Handler

type Router interface {
  // Middlewares
  Use(mw ...Middleware)

  // Composition
  Group(prefix string, fn func(Router))     // sub-tree under prefix
  GroupWith(mw Middleware, fn func(Router)) // sub-tree with middleware applied
  Mount(prefix string, h http.Handler)      // attach a ready handler (e.g., GraphQL)

  // HTTP handlers
  Handle(method, path string, h http.Handler)
  GET(path string, h http.Handler)
  POST(path string, h http.Handler)
  PUT(path string, h http.Handler)
  DELETE(path string, h http.Handler)

  // Defaults / integration
  SetNotFound(h http.Handler)
  SetMethodNotAllowed(h http.Handler)
  SetError(func(http.ResponseWriter, *http.Request, error))

  ServeHTTP(http.ResponseWriter, *http.Request)
}
```

### Key points

* `Middleware` is plain `func(http.Handler) http.Handler)` — easy to reuse and test.
* `GroupWith` scopes middleware to a **subtree** (ideal for auth-protected areas).
* `Mount` attaches a fully built `http.Handler` (used to mount GraphQL under `cfg.ServerGraphql.Path`).
* `SetError` is an optional platform sink for surfacing errors (paired with recovery).

## How contexts use it

Each context exports a registrar and never imports a concrete router:

```go
// internal/user/adapter/primary/http/handler/register.go
func RegisterHTTP(r ports.Router, h *Handler) {
  r.Group("/v1/users", func(ur ports.Router) {
    ur.POST("/create", http.HandlerFunc(h.Create)) // public
  })
}

// internal/user/adapter/primary/http/handler/register_protected.go
func RegisterHTTPProtected(r ports.Router, h *Handler, mw ports.Middleware) {
  r.Group("/v1/users", func(ur ports.Router) {
    ur.GroupWith(mw, func(pr ports.Router) {
      pr.PUT("/", http.HandlerFunc(h.UpdateUser))
      pr.PUT("/password", http.HandlerFunc(h.UpdateUserPassword))
      pr.DELETE("/", http.HandlerFunc(h.SoftDeleteUser))
    })
  })
}
```

## Swapping the router

* Implement `ports.Router` under `internal/platform/server/http/router/<impl>/<impl>_router.go`.
* Update the platform composer to instantiate your new router instead of chi.
* No changes are required in context registrars.

## Best practices

* **Do not** import `chi/gin/echo` in context packages; use only `ports.Router`.
* Apply **global** platform middlewares in the HTTP composer; apply **domain** middlewares with `GroupWith` in the context registrar.
* Mount GraphQL via `Mount(...)` instead of re-declaring GET/POST for it.
* Keep `SetNotFound`, `SetMethodNotAllowed`, and `SetError` centralized in the platform to ensure consistent responses/telemetry.
