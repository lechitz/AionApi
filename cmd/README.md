# Application Entrypoints (`cmd`)

**Path:** `cmd`

## Overview

This folder contains application binaries.
Current production entrypoint is `cmd/api`.

## Structure

| Folder | Role |
| --- | --- |
| `api/` | Main production server bootstrap |

## Design Notes

- Entrypoints should orchestrate and delegate, not implement domain logic.
- Keep startup concerns centralized (config, lifecycle, observability).
- Development tooling belongs under `hack/`.

## Package Improvements

- Add a compact boot sequence section linking to platform startup docs.
- Document expected env profile for each entrypoint.
- Add explicit convention for naming new entrypoint folders.
- Add link to release/build pipeline docs for binaries.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../README.md)
<!-- doc-nav:end -->
