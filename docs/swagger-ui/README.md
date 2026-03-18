# Swagger UI Static Bundle

**Path:** `docs/swagger-ui`

## Purpose

This folder contains the static Swagger UI bundle used for published documentation hosting.

It is not the runtime HTTP implementation.
The live API server mounts Swagger with `http-swagger` inside `internal/platform/server/http/composer.go`.

## Current Relationship To Runtime

- static hosting: files in this folder
- live runtime mount: `${HTTP_CONTEXT}${HTTP_SWAGGER_MOUNT_PATH}`
- docs alias redirect: `${HTTP_CONTEXT}${HTTP_DOCS_ALIAS_PATH}` -> `${HTTP_CONTEXT}${HTTP_SWAGGER_MOUNT_PATH}/index.html`
- contract source: `contracts/openapi/swagger.json`

## Files

| File | Purpose |
| --- | --- |
| `index.html` | Static Swagger UI entrypoint |
| `swagger-ui-bundle.js`, `swagger-ui.css` | Vendored upstream Swagger UI assets |
| `custom.css` | Local presentation overrides for hosted docs |

## Rules

- Regenerate the OpenAPI contract before reviewing or publishing UI output.
- Keep static asset paths compatible with the GitHub Pages/docs hosting layout.
- Do not treat this folder as the source of truth for REST behavior; it renders the generated contract.
- After asset upgrades, verify both hosted docs and the runtime Swagger mount.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../index.md)
<!-- doc-nav:end -->
