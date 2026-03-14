# Docker Environment Template

**Path:** `infrastructure/docker/environments/example`

## Overview

This folder is the baseline template for creating new Docker environment profiles.
It defines a complete `.env.example` with all required application, infrastructure, and observability variables.

## Package Scope

| Area | Responsibility |
| --- | --- |
| Variable baseline | Provide the canonical variable set required by the stack |
| Onboarding support | Offer safe placeholders for local/profile setup |
| Consistency guard | Reduce drift when creating new environment folders |

## Files

| File | Purpose |
| --- | --- |
| `.env.example` | Full environment variable template with grouped sections and guidance |
| `README.md` | How to use this template when creating a new environment profile |

## Variable Groups in `.env.example`

| Group | Examples |
| --- | --- |
| General app metadata | `APP_NAME`, `APP_ENV`, `APP_VERSION` |
| OpenTelemetry | `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_SERVICE_NAME` |
| HTTP server | `HTTP_HOST`, `HTTP_PORT`, `HTTP_CONTEXT`, `HTTP_API_ROOT` |
| GraphQL server | `GRAPHQL_PATH`, timeout and limits |
| Database | `DB_*`, pooling and retry settings |
| Redis/cache | `CACHE_*` |
| Auth/JWT | `SECRET_KEY` |

## How to Create a New Environment

1. Copy this template into your target profile folder (`dev`, `staging`, etc.).
2. Replace placeholders (`<your_...>`) with profile-specific values.
3. Keep naming and variable keys aligned with application config structs.
4. Validate by running the corresponding compose/make target.

Example:

```bash
cp infrastructure/docker/environments/example/.env.example infrastructure/docker/environments/dev/.env.dev
```

## Design Notes

- This folder should contain only safe templates, never real secrets.
- Environment-specific runtime behavior belongs to profile folders under `infrastructure/docker/environments/`.
- The template is a contract aid between app config and infrastructure wiring.

## Package Improvements

- Add a validation script/target that checks `.env` completeness against required config fields.
- Document mandatory vs optional variables explicitly in `.env.example` comments.
- Consider splitting template sections by service ownership (API, DB, cache, observability) with short examples.
- Add a tiny “minimal local setup” sample alongside full template for faster first run.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../README.md)
<!-- doc-nav:end -->
