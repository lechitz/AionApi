# Platform HTTP Router Port

**Path:** `internal/platform/server/http/ports`

## Overview

This package defines the HTTP routing contract used by platform and context adapters.
It is the abstraction that keeps contexts independent from concrete router frameworks (chi, gin, echo, etc.).

## Package Scope

| Area | Responsibility |
| --- | --- |
| Router abstraction | Define the `Router` interface used by HTTP registrars |
| Middleware abstraction | Define framework-neutral middleware signature |
| Composition contract | Standardize grouping, mounting, and fallback handler registration |

## Files

| File | Purpose |
| --- | --- |
| `router.go` | Declares `Middleware` and `Router` interfaces |

## Core Contracts

| Contract | Description |
| --- | --- |
| `type Middleware func(http.Handler) http.Handler` | Shared middleware signature across the HTTP layer |
| `type Router interface { ... }` | Route registration, grouping, mounting, fallbacks, and `ServeHTTP` |

## Router Capabilities

| Capability | Methods |
| --- | --- |
| Global middleware | `Use(...)` |
| Route tree composition | `Group(...)`, `GroupWith(...)` |
| Handler mounting | `Mount(prefix, handler)` |
| HTTP method mapping | `Handle`, `GET`, `POST`, `PUT`, `DELETE` |
| Fallback behavior | `SetNotFound`, `SetMethodNotAllowed`, `SetError` |
| Entrypoint | `ServeHTTP` |

## Usage Pattern

```go
func RegisterHTTP(r ports.Router, h *Handler) {
    r.Group("/v1/users", func(gr ports.Router) {
        gr.POST("/create", http.HandlerFunc(h.Create))
    })
}
```

## Design Notes

- Context packages should depend on this port only, never on concrete router libraries.
- Router implementation details belong to `internal/platform/server/http/router/...`.
- Fallback handlers should be centralized at composer level for consistent transport behavior.

## Package Improvements

- Add contract tests for any router adapter implementation to ensure full parity with `Router` behavior.
- Consider adding short semantic comments in `router.go` for `SetError` execution expectations.
- Evaluate whether `PATCH` support is needed in the port for future API evolution.
- Add a small reference in platform docs to clarify when to use `Mount` vs `Group` + method handlers.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../README.md)
<!-- doc-nav:end -->
