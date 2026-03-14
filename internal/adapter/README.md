# Shared Adapter Layer

**Path:** `internal/adapter`

## Overview

Shared adapter layer containing transport-facing primary adapters and infrastructure-facing secondary adapters.
It provides cross-context adapter infrastructure and conventions.

## Subpackages

| Subpackage | Role |
| --- | --- |
| `primary/` | Inbound transport adapters (client -> usecase boundary) |
| `secondary/` | Outbound infra adapters (usecase -> external systems) |

## Design Notes

- Keep adapter responsibilities directional and explicit.
- Keep core boundaries clean: no domain rules in adapters.
- Detailed behavior for each side is documented in respective subpackage READMEs.

## Package Improvements

- Add adapter boundary checklist (inbound vs outbound responsibilities).
- Add sample request flow bridging primary -> core -> secondary.
- Add lint/convention checks for forbidden dependency direction.
- Add adapter testing guidance by adapter type.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
