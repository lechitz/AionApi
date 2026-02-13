# Swagger UI Static Bundle

**Path:** `docs/swagger-ui`

## Overview

Static Swagger UI assets used to browse the generated OpenAPI contract.
This folder is intended for documentation hosting, not runtime API logic.

## Files

| File | Purpose |
| --- | --- |
| `index.html` | Swagger UI entrypoint |
| `swagger-ui-bundle.js`, `swagger-ui.css` | Upstream Swagger UI distribution assets |
| `custom.css` | Local visual overrides |

## Usage Notes

- Regenerate OpenAPI before reviewing UI output.
- Keep relative references compatible with static hosting paths.
- Validate UI after bumping Swagger UI distribution files.

## Package Improvements

- Add source/version note for bundled Swagger UI release.
- Add quick test checklist after asset upgrades.
- Add minimal customization policy to keep overrides maintainable.
- Add a link to OpenAPI generation workflow (`make swag`).

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../index.md)
<!-- doc-nav:end -->
