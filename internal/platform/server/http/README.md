# Platform HTTP Server Layer

**Path:** `internal/platform/server/http`

## Overview

This package composes the application HTTP surface: router adapter, middlewares, generic handlers, REST registrations, GraphQL mount, Swagger mount, and health endpoints.
It is the platform entrypoint for transport-level HTTP orchestration.

## Package Scope

| Area | Responsibility |
| --- | --- |
| Handler composition | Build the main HTTP handler tree (`ComposeHandler`) |
| Server configuration | Build `http.Server` with config-driven timeouts and limits |
| Platform wiring | Apply cross-cutting middlewares and fallback handlers |
| Endpoint mounting | Mount REST adapters, GraphQL, Swagger, and health routes |
| HTTP sentinel errors | Provide shared transport errors (`404`, `405`, `500`) |

## Subpackages

| Subpackage | Role |
| --- | --- |
| `ports/` | Framework-agnostic router contract |
| `router/` | Concrete router adapter implementations (currently `chi`) |
| `middleware/` | Cross-cutting HTTP middleware chain |
| `generic/` | Generic handlers (`health`, `not found`, `method not allowed`, `recovery/error`) |
| `utils/` | Shared response/error/cookie helpers |
| `errors/` | Shared sentinel HTTP errors |

## Core Files

| File | Purpose |
| --- | --- |
| `composer.go` | Assembles full handler graph (middlewares + routes + GraphQL + Swagger + health) |
| `server.go` | Builds `http.Server` from config and composed handler |
| `http_constants.go` | Constants for routing/mount defaults and logging |

## Composition Flow (`ComposeHandler`)

1. Instantiate router adapter (`chi.New()`).
2. Register global middlewares (`requestid`, `recovery`, `cors`).
3. Set fallback handlers (`NotFound`, `MethodNotAllowed`, `Error`).
4. Resolve mount points from config (`context`, `swagger`, `docs`, `health`).
5. Mount Swagger UI and docs alias under API context.
6. Mount REST modules under API root.
7. Build and mount GraphQL handler; wrap with `servicetoken` middleware.
8. Wrap main router with `otelhttp` instrumentation.
9. Expose health routes via a dedicated mux path (including backward-compatible path).

## Server Build (`Build` / `FromHTTP`)

| Concern | Source |
| --- | --- |
| Listen address | `cfg.ServerHTTP.Host` + `cfg.ServerHTTP.Port` |
| Timeouts/limits | `cfg.ServerHTTP.*Timeout`, `MaxHeaderBytes` |
| Base context | Application context propagated via `BaseContext` |

## Design Notes

- This layer must remain transport/platform-focused; no domain usecase logic.
- Context modules register endpoints through adapters; composition happens here.
- Subpackage READMEs contain detailed contracts/behavior; this README is the integration-level view.

## Package Improvements

- Add integration tests for `ComposeHandler` route map (Swagger, docs alias, GraphQL mount, health aliases).
- Ensure middleware ordering in code matches documented recommendation (`recovery` outermost) or update docs/code for consistency.
- Consider adding a small table mapping mounted routes to owning adapters for discoverability.
- Add a short section documenting failure behavior when GraphQL handler composition fails inside route setup.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../README.md)
<!-- doc-nav:end -->
