# Shared Constants

**Path:** `internal/shared/constants`

## Overview

Centralized string and key constants for headers, claims, context keys, logging, and tracing attributes.
This package removes magic strings and keeps cross-cutting telemetry consistent.

## Responsibilities

| Area | Responsibility |
| --- | --- |
| Key centralization | Single source of truth for common keys |
| Consistency | Stabilize logs/headers/traces across contexts |
| Safety | Reduce typo and drift risk across packages |

## Design Notes

- Keep only constants and minimal helper types.
- Do not place business/domain logic here.
- Preserve key stability to avoid dashboard/query regressions.

## Package Improvements

- Add a key taxonomy map by domain concern.
- Add automated duplicate-key detection.
- Add changelog guidance for key renames/deprecations.
- Add examples for when to use `ctxkeys` vs `commonkeys`.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
