# Docker Infrastructure

**Path:** `infrastructure/docker`

## Purpose

This folder owns the container image and compose wiring used to run `AionApi` in local, personal, and prod-like profiles.

## Current Layout

| Path | Responsibility |
| --- | --- |
| `Dockerfile` | multi-stage image that builds `aion-api` and `aion-api-outbox-publisher` |
| `scripts/entrypoint.sh` | default container entrypoint; starts `aion-api` |
| `environments/dev/` | integrated hot-reload stack used in the multi-repo workspace |
| `environments/my/` | personal isolated stack with its own compose and env files |
| `environments/prod/` | prod-like compose profile |
| `environments/example/` | env template for creating new profiles |

## Operational Flows

```bash
make build-dev
make dev
make rebuild-dev
make my
make prod-up
```

The `dev` profile also runs the outbox publisher, observability stack, and sibling services from `aion-chat`, `aionapi-dashboard`, `aion-ingest`, and `aion-streams`.

## Boundaries

- an isolated clone of `AionApi` can build the image, but the full `make dev` flow assumes the complete `/Aion` workspace
- keep environment-specific values in profile env files, not in the root `Dockerfile`
- if container behavior differs from runtime code, the entrypoint script, compose profile, and Make targets are the canonical sources
- this folder owns image and compose wiring only; app configuration remains in `internal/platform/config`

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
