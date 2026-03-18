# Platform HTTP Router Port

**Path:** `internal/platform/server/http/ports`

## Purpose

This package defines the framework-neutral routing contract consumed by HTTP registrars and the platform composer.

## Current Contract

| Capability | Methods |
| --- | --- |
| Middleware | `Use(...)` |
| Tree composition | `Group(...)`, `GroupWith(...)`, `Mount(...)` |
| Method mapping | `Handle`, `GET`, `POST`, `PUT`, `DELETE` |
| Fallbacks | `SetNotFound`, `SetMethodNotAllowed`, `SetError` |
| Entrypoint | `ServeHTTP` |

The concrete implementation lives under `internal/platform/server/http/router/chi`.

## Usage Pattern

```go
func RegisterHTTP(r ports.Router, h *Handler) {
	r.Group("/v1/users", func(gr ports.Router) {
		gr.POST("/create", http.HandlerFunc(h.Create))
	})
}
```

## Boundaries

- context packages should depend on this port only, never on router-specific APIs
- fallback behavior remains centralized in the platform composer for consistent transport semantics
- if this contract changes, the chi adapter and composed server tests must change with it

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../README.md)
<!-- doc-nav:end -->
