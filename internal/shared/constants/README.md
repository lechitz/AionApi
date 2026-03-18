# Shared Constants

**Path:** `internal/shared/constants`

## Purpose

This package centralizes the stable cross-cutting keys reused by multiple parts of the runtime.

## Current Namespaces

| Namespace | Role |
| --- | --- |
| `claimskeys/` | JWT claim names |
| `commonkeys/` | log fields, request keys, cookie names, and other shared labels |
| `ctxkeys/` | typed context keys |
| `roles/` | shared role names |
| `tracingkeys/` | legacy HTTP/request tracing attributes kept for compatibility |

## Boundaries

- keep only constants and minimal helper types here
- tracing and status strings specific to one bounded context should stay local to that context
- renaming shared keys is a compatibility change because it can affect auth, logs, traces, dashboards, and clients at once

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
