# Router Port (Platform -> HTTP)

**Path:** `internal/platform/server/http/ports`

## Purpose and Main Capabilities

- Define the framework-agnostic routing contract (`Router`).
- Keep contexts decoupled from concrete routers.
- Standardize how routes and middleware are registered.

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

- `Middleware` is `func(http.Handler) http.Handler` (easy to reuse and test).
- `GroupWith` scopes middleware to a subtree (ideal for auth-protected areas).
- `Mount` attaches a fully built handler (GraphQL, etc.).
- `SetError` is used by recovery/error flows (not called by the router itself).

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

## Recommended Practices Visible Here

- Do not import `chi/gin/echo` in contexts; use only `ports.Router`.
- Apply global middlewares in the composer; domain middlewares via `GroupWith`.
- Mount GraphQL via `Mount(...)` instead of per-method handlers.
- Keep NotFound/MethodNotAllowed/Error centralized for consistent responses.

## What Should NOT Live Here

- Concrete router implementations.
- Domain logic or handler code.
