# Infrastructure Layer

**Path:** `infrastructure`

## Overview

Operational infrastructure assets for container runtime, database lifecycle, and observability stack.
This layer supports local development and deployment-oriented workflows.

## Subpackages

| Subpackage | Responsibility |
| --- | --- |
| `docker/` | Image and environment runtime definitions |
| `db/` | Migration and seed SQL assets |
| `observability/` | Logs/metrics/traces infrastructure configs |

## Design Notes

- Keep infrastructure concerns separate from application/business code.
- Favor reproducible, versioned configuration assets.
- Keep sensitive runtime values outside committed docs/config templates.

## Package Improvements

- Add infra quick-start flowchart (docker + db + observability bootstrap).
- Add “who owns what” matrix per infrastructure subpackage.
- Add policy section for secret management across profiles.
- Add compatibility matrix between app version and infra dependencies.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../README.md)
<!-- doc-nav:end -->
