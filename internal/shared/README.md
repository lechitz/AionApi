# Shared Cross-Cutting Layer

**Path:** `internal/shared`

## Purpose

`internal/shared` holds the small set of cross-cutting assets that multiple boundaries rely on.

## Current Area

| Area | Role |
| --- | --- |
| `constants/` | shared keys for claims, headers, context values, roles, and selected tracing attributes |

## Boundaries

- this layer must stay minimal and stable
- do not move business logic or context-specific semantics here
- prefer local constants in the owning package unless the value is intentionally shared across boundaries

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
