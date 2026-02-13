# Platform HTTP Utilities Layer

**Path:** `internal/platform/server/http/utils`

## Overview

This package aggregates shared HTTP utility modules used by platform and context adapters.
It provides a consistent foundation for response writing, semantic error handling, and auth cookie operations.

## Subpackages

| Subpackage | Responsibility |
| --- | --- |
| `sharederrors/` | HTTP-facing semantic errors and error-to-status mapping |
| `httpresponse/` | Standard response envelope and response writer helpers |
| `cookies/` | Auth/refresh cookie lifecycle and extraction helpers |

## Layer Responsibilities

| Concern | Utility module |
| --- | --- |
| Standard JSON success/error envelope | `httpresponse` |
| Stable semantic error taxonomy | `sharederrors` |
| HTTP status mapping from semantic errors | `sharederrors` + `httpresponse` |
| Cookie security and token transport | `cookies` |

## Typical Request/Response Path

1. Handler/usecase returns success payload or semantic error.
2. `sharederrors` defines the semantic category.
3. `httpresponse` maps status and writes standardized JSON response.
4. `cookies` helpers set/clear/extract auth cookies when needed.

## Design Notes

- Keep all utilities in this layer transport-focused and reusable across contexts.
- Domain/business rules must not be added to utility packages.
- Subpackage READMEs contain implementation details; this README is the high-level integration view.

## Package Improvements

- Add cross-package tests validating end-to-end behavior (`sharederrors` -> `httpresponse`) for common error classes.
- Define a shared guideline for when to include `details` in error responses across environments.
- Align cookie token names and other transport keys with shared constants to remove hardcoded literals.
- Add a “recommended usage matrix” mapping common handler scenarios to utility helpers.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../../README.md)
<!-- doc-nav:end -->
