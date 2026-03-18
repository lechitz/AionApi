# Internal Application Layer

**Path:** `internal`

## Purpose

`internal` contains the application code that is not meant to be imported outside the module.
The repo is organized around bounded contexts plus a small set of cross-cutting runtime layers.

## Current Areas

| Area | Role |
| --- | --- |
| `admin`, `audit`, `auth`, `category`, `chat`, `eventoutbox`, `realtime`, `record`, `tag`, `user` | bounded contexts with core ports, use cases, and local adapters |
| `adapter/` | shared adapter infrastructure reused across contexts |
| `platform/` | config, Fx wiring, runtime services, ports, and server composition |
| `shared/` | stable cross-cutting constants and key namespaces |

## Boundaries

- context business rules belong inside the owning bounded context
- shared transport or infra helpers belong in `adapter/` or `platform/` only when they are reused across contexts
- `shared/` stays intentionally small; context-specific constants should remain local when possible

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../README.md)
<!-- doc-nav:end -->
