# Database Infrastructure

**Path:** `infrastructure/db`

## Overview

Database infrastructure assets for schema lifecycle and seed datasets.
This package is split into migration contracts and seed scripts.

## Subpackages

| Subpackage | Responsibility |
| --- | --- |
| `migrations/` | Versioned schema evolution |
| `seed/` | Deterministic local/test data provisioning |

## Design Notes

- Keep schema changes migration-driven.
- Keep seed data separate from schema evolution.
- Align SQL artifacts with repository expectations in DB adapters.

## Package Improvements

- Add a simple compatibility matrix (migration version vs seed assumptions).
- Add DB bootstrap script for first-time local setup.
- Add schema verification query snippets for quick checks.
- Add naming conventions for new SQL artifacts.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
