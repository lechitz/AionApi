# API Collections

**Path:** `docs/collections`

## Overview

This folder contains API client collections (Postman/Insomnia) used for manual QA and local integration checks.
OpenAPI remains the source contract; collections are consumer-friendly artifacts.

## Scope

| Area | Responsibility |
| --- | --- |
| Manual testing support | Provide ready-to-run request collections |
| Contract consumption | Reflect REST contract for external tooling |

## Design Notes

- Keep environment placeholders (`{{baseURL}}`, tokens) in collection files.
- Never commit real secrets/tokens.
- Update collections when contract changes materially affect consumers.

## Package Improvements

- Add collection versioning notes linked to API releases.
- Add validation step to ensure collections import cleanly.
- Add quick start section with expected local environment variables.
- Add mapping table from collection folders to API domains.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
