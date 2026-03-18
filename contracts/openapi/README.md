# OpenAPI Contract

**Path:** `contracts/openapi`

## Purpose

This folder publishes the consumer-facing REST contract used by Swagger UI, static docs hosting, and client generation.

Authority order for REST remains:

1. live handler behavior and annotated routes in `internal/*/adapter/primary/http/handler`
2. generated OpenAPI artifacts in this folder
3. narrative docs such as `docs/swagger-ui/README.md`

## Current Files

| File | Purpose |
| --- | --- |
| `swagger.yaml` | Published OpenAPI document in YAML |
| `swagger.json` | Published OpenAPI document in JSON |
| `docs.go` | `swag`-generated Go metadata consumed by Swagger tooling |

## Current Surface

The published REST contract documents the annotated HTTP endpoints owned by the runtime composer, including:

- auth session flows
- user and admin management endpoints
- chat endpoints
- audit diagnostics endpoints

GraphQL is intentionally out of scope for this folder.

## Workflow

Regenerate the contract after changing annotated HTTP handlers:

```bash
make swag
```

Common downstream generation still starts from `swagger.json`:

```bash
openapi-generator-cli generate -i contracts/openapi/swagger.json -g typescript-axios -o clients/typescript
openapi-generator-cli generate -i contracts/openapi/swagger.json -g python -o clients/python
```

## Rules

- Handler behavior and annotations must change in the same PR as contract regeneration.
- Treat `swagger.yaml`, `swagger.json`, and `docs.go` as generated artifacts.
- Validate the published UI at `${HTTP_CONTEXT}${HTTP_SWAGGER_MOUNT_PATH}` or `${HTTP_CONTEXT}${HTTP_DOCS_ALIAS_PATH}` after regeneration.
- Do not describe GraphQL, Kafka, or SSE contracts here; keep this folder REST-only.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
