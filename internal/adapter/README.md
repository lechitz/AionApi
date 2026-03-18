# Shared Adapter Layer

**Path:** `internal/adapter`

## Purpose

`internal/adapter` holds adapter code that is shared across more than one bounded context.

## Current Split

| Subpackage | Role |
| --- | --- |
| `primary/` | shared inbound transport infrastructure |
| `secondary/` | shared outbound infrastructure implementations |

## Boundaries

- context-specific HTTP, DB, cache, or storage adapters should stay inside the owning context
- only cross-context adapter infrastructure belongs here
- adapters translate transport or vendor behavior; they do not own business orchestration

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
