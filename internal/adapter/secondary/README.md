# Secondary Adapters Layer

**Path:** `internal/adapter/secondary`

## Overview

Secondary adapters implement output ports and connect core usecases to external systems.
They are responsible for IO translation and infrastructure error adaptation.

## Responsibilities

| Area | Responsibility |
| --- | --- |
| Port implementation | Implement output contracts from core/platform ports |
| Infra translation | Convert IO/vendor behavior into semantic application behavior |
| Boundary observability | Emit logs/traces around external calls |

## Design Notes

- Keep infra-specific details out of core.
- Keep adapters replaceable and minimal.
- Do not put business orchestration in secondary adapters.

## Package Improvements

- Add implementation matrix (port -> concrete adapter package).
- Add error translation guidelines by adapter type.
- Add retry/timeout policy references for external IO.
- Add adapter conformance testing guide.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
