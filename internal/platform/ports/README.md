# Platform Output Ports

**Path:** `internal/platform/ports`

## Overview

Cross-cutting output port contracts (cache, logger, hasher, keygen, etc.) used across contexts.
These interfaces keep core/business logic decoupled from vendor implementations.

## Responsibilities

| Area | Responsibility |
| --- | --- |
| Contract definition | Define minimal technology-agnostic interfaces |
| Decoupling | Prevent direct dependency on infra libraries in core |
| Testability | Enable mock generation for unit tests |

## Design Notes

- Keep interfaces focused and small.
- Avoid vendor/framework types in port contracts.
- Treat ports as stable contracts used by multiple contexts.

## Package Improvements

- Add contract index table linking each port to adapters.
- Add versioning policy for breaking contract changes.
- Add examples showing preferred usage from usecases.
- Add CI check for stale/unmocked new ports.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
