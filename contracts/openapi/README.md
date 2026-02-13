# OpenAPI Contract

**Path:** `contracts/openapi`

## Overview

This folder stores the REST API OpenAPI specification used by Swagger UI and client generation.
It is the contract artifact for HTTP consumers.

## Files

| File | Purpose |
| --- | --- |
| `swagger.yaml` | OpenAPI document (YAML) |
| `swagger.json` | OpenAPI document (JSON) |
| `docs.go` | swag-generated Go metadata |

## Common Workflows

### Regenerate spec

```bash
make swag
```

### Generate clients

```bash
openapi-generator-cli generate -i contracts/openapi/swagger.json -g typescript-axios -o clients/typescript
openapi-generator-cli generate -i contracts/openapi/swagger.json -g python -o clients/python
```

## Design Notes

- Treat this folder as generated contract output from annotated handlers.
- Keep REST contract updates synchronized with endpoint behavior changes.
- Avoid manual edits on generated artifacts when regeneration is available.

## Package Improvements

- Add CI check to detect drift between annotations and committed spec files.
- Add changelog section for breaking REST contract changes.
- Add contract linting step (spectral/openapi-cli) in verify workflow.
- Add a short “auth scheme + base path” quick reference table.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
