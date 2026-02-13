# Fx Application Wiring

**Path:** `internal/platform/fxapp`

## Overview

Dependency graph composition and lifecycle wiring using Uber Fx.
This package assembles infrastructure, platform, and domain dependencies for runtime startup.

## Responsibilities

| Area | Responsibility |
| --- | --- |
| Module composition | Group providers/invokes by platform area |
| Dependency injection | Build runtime object graph |
| Lifecycle hooks | Coordinate start/stop and graceful shutdown |

## Design Notes

- Keep module boundaries explicit.
- Prefer provider granularity that matches bounded contexts.
- Log and surface startup/shutdown failures clearly.

## Package Improvements

- Add module dependency graph diagram.
- Add startup failure troubleshooting section.
- Add conventions for adding new providers/invokes.
- Add smoke test for app boot sequence.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
