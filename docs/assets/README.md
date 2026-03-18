# Documentation Assets

**Path:** `docs/assets`

## Purpose

This folder holds static assets consumed by the MkDocs site.
It does not ship runtime assets for the API or dashboard.

## Current Files

| Path | Used by | Purpose |
| --- | --- | --- |
| `logo/Aion.png` | `mkdocs.yml` logo and favicon | branding for the published docs site |
| `stylesheets/extra.css` | `mkdocs.yml` `extra_css` | small visual overrides for the docs theme |

## Boundaries

- keep file names stable when they are referenced from `mkdocs.yml`
- store only documentation-facing assets here
- runtime media for product surfaces belongs in the owning repo or delivery bucket, not under `docs/assets`

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../index.md)
<!-- doc-nav:end -->
