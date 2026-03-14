# Shared Cross-Cutting Layer

**Path:** `internal/shared`

## Overview

Cross-context shared assets that do not contain business logic.
Currently focused on shared constants and key namespaces used across adapters/usecases/platform.

## Subpackages

| Subpackage | Role |
| --- | --- |
| `constants/` | Claims/header/context/log/tracing key definitions |

## Design Notes

- Keep this layer minimal and stable.
- Avoid introducing context business rules here.
- Use shared keys to reduce duplication and observability drift.

## Package Improvements

- Add deprecation policy for renamed shared keys.
- Add guidance for introducing new constant namespaces.
- Add references to dashboards/log queries impacted by key changes.
- Add integrity checks for duplicate/overlapping constants.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
