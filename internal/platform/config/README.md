# Platform Configuration

**Path:** `internal/platform/config`

## Overview

Typed configuration loading and validation for platform/domain runtime.
It centralizes env parsing, normalization, and fail-fast checks.

## Responsibilities

| Area | Responsibility |
| --- | --- |
| Env loading | Parse env vars into typed config structs |
| Validation | Fail early on invalid/missing required values |
| Normalization | Apply safe defaults/path normalization/time constraints |
| Secret bootstrap | Support generated temporary secrets for local usage |

## Design Notes

- Configuration should be explicit and typed.
- Keep validation close to config model.
- Avoid hidden defaults that can mask misconfiguration.

## Package Improvements

- Add env key reference table generated from config tags.
- Add test coverage for edge-case normalization rules.
- Add explicit production safety notes for secret handling.
- Add lint/check to detect undocumented new env keys.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
