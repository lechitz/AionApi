# Platform HTTP Router Adapters

**Path:** `internal/platform/server/http/router`

## Overview

This package contains concrete router adapters that implement the `ports.Router` contract.
It isolates framework-specific routing details from context modules and keeps route registration portable.

## Package Scope

| Area | Responsibility |
| --- | --- |
| Port implementation | Implement `internal/platform/server/http/ports.Router` |
| Framework encapsulation | Hide chi-specific APIs from contexts and registrars |
| Middleware composition | Apply global and scoped middlewares through router adapter methods |
| Fallback wiring | Configure `NotFound`, `MethodNotAllowed`, and optional error callback hooks |

## Subpackages

| Subpackage | Role |
| --- | --- |
| `chi/` | Current `chi/v5` adapter implementing the router port |

## Current Adapter Summary (`chi/`)

| Capability | Implementation notes |
| --- | --- |
| Global middleware | `Use(...)` with nil-guard |
| Route groups | `Group(prefix, fn)` and `GroupWith(mw, fn)` |
| Mounted handlers | `Mount(prefix, handler)` |
| Method routing | `Handle`, `GET`, `POST`, `PUT`, `DELETE` |
| Fallback handlers | `SetNotFound`, `SetMethodNotAllowed` |
| Error callback | `SetError` stores callback (not automatically invoked by chi) |

## Design Notes

- Context registrars should depend only on `ports.Router`, never on `chi` directly.
- Adapter remains thin: route wiring only, no business or transport mapping logic.
- Router swapping should be done in composer wiring, not in domain/context packages.

## Package Improvements

- Add contract tests validating that `chi` adapter behavior matches `ports.Router` expectations.
- Clarify and/or remove `SetError` if it is not consumed by platform flow to avoid dead API surface.
- Consider adding optional support for additional verbs (e.g., `PATCH`) if required by upcoming endpoints.
- Add a short “how to add a new adapter” section in this README with required method parity checklist.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../README.md)
<!-- doc-nav:end -->
