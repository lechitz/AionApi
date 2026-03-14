# Platform Layer

**Path:** `internal/platform`

## Overview

Cross-cutting platform layer for configuration, dependency graph assembly, observability, shared ports, and server runtime composition.
It provides the operational foundation used by all bounded contexts.

## Subpackages

| Subpackage | Role |
| --- | --- |
| `config/` | Typed env loading and validation |
| `fxapp/` | Dependency graph and lifecycle wiring |
| `observability/` | Tracing/metrics provider bootstrap |
| `ports/` | Shared output port contracts |
| `httpclient/` | Shared instrumented outbound HTTP client |
| `server/` | HTTP server composition and transport wiring |

## Design Notes

- Keep platform code domain-agnostic.
- Keep runtime wiring centralized and explicit.
- Keep subpackages independently documented for deeper details.

## Package Improvements

- Add platform startup sequence diagram linking all subpackages.
- Add compatibility table (config keys -> platform modules).
- Add observability integration checklist for new modules.
- Add convention guide for introducing new platform services.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
