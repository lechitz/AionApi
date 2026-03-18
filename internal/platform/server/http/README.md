# Platform HTTP Server Layer

**Path:** `internal/platform/server/http`

## Purpose

This package composes the complete HTTP transport surface for `AionApi`.
It owns router creation, middleware wiring, Swagger/docs mounts, REST registrar composition, GraphQL mount, and health endpoints.

## Current Composition

| Concern | Current behavior |
| --- | --- |
| Main router | `chi` adapter created in `ComposeHandler` |
| Global middleware order | `requestid` -> `recovery` -> `cors` |
| Fallback handlers | `NotFound`, `MethodNotAllowed`, and generic error callback are wired from `generic/handler` |
| Swagger/docs | mounted under `cfg.ServerHTTP.Context` with alias redirect to `.../swagger/index.html` |
| REST routes | mounted under `cfg.ServerHTTP.Context + cfg.ServerHTTP.APIRoot` |
| GraphQL | mounted inside the API root at `cfg.ServerGraphql.Path`, wrapped by `servicetoken` |
| OTel HTTP wrapper | wraps the main router, not the dedicated health mux |
| Health routes | exposed separately at both `${context}${health}` and `${context}${apiRoot}${health}` |

## Conditional Route Registration

`registerDomainRoutes` mounts REST modules only when their dependencies are present:

- auth
- user
- admin
- chat
- audit
- realtime

GraphQL handler construction is also dependency-driven and mounted only after successful composition.

## Health Exception

Health endpoints intentionally bypass the instrumented main router.
They currently run through:

- `requestid`
- `cors`

They do not go through `otelhttp` or the main router fallback chain.

## Key Files

| File | Purpose |
| --- | --- |
| `composer.go` | route/middleware composition and mount logic |
| `server.go` | `http.Server` construction from config |
| `http_constants.go` | default route and mount constants |

## Boundaries

- No domain usecase logic belongs here.
- Bounded contexts register handlers through adapters; this package only composes them.
- Shared transport primitives live in `middleware/`, `generic/`, `router/`, `ports/`, and `utils/`.

## Validate

```bash
go test ./internal/platform/server/http/...
make verify
```

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../README.md)
<!-- doc-nav:end -->
