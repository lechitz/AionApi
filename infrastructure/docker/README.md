# infrastructure/docker

Docker assets for building and running the AionAPI stack. This folder defines how images and compose profiles are structured.

## Package Composition

- `Dockerfile`
  - Multi-stage build for the API binary.
- `environments/`
  - Compose overlays and env files for dev/prod/example.
- `scripts/`
  - Container entrypoints and helper scripts.

## Flow (Where it comes from -> Where it goes)

Dockerfile + env profiles -> docker compose -> running services

## Why It Was Designed This Way

- Keep runtime profiles isolated (dev vs prod).
- Keep Docker concerns out of application code.
- Use reproducible builds and consistent startup.

## Recommended Practices Visible Here

- Multi-stage build for small images.
- Per-environment env files and compose overlays.
- Entry-point scripts are simple and explicit.

## Common Commands

```bash
make build-dev
make dev-up
make dev-down
make clean-dev
make build-prod
make prod-up
```

## What Should NOT Live Here

- Business logic or app behavior.
- Real secrets committed to git.
- Environment-specific overrides in shared templates.
