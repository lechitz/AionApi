# infrastructure/docker

Docker assets for building and running the AionAPI stack. This folder defines how images and compose profiles are structured.

## Package Composition

- `Dockerfile`
  - Multi-stage build for the API binary (production).
- `environments/`
  - `dev/`
    - `Dockerfile.dev` - Development image with Air hot reload
    - `docker-compose-dev.yaml` - Dev orchestration
    - `.env.dev` - Dev environment variables
  - `prod/`
    - `docker-compose-prod.yaml` - Production orchestration
    - `.env.prod` - Production environment variables
  - `example/`
    - `.env.example` - Template for environment variables
- `scripts/`
  - Container entrypoints and helper scripts.

## Development vs Production

**Development (Dockerfile.dev):**
- Uses full `golang:1.25` image (includes build tools)
- Installs Air for hot reload
- Source code mounted via volumes (not copied)
- Fast iteration: changes detected in 3-5 seconds

**Production (Dockerfile):**
- Multi-stage build with Alpine (minimal size)
- Compiled binary only (no source code)
- No hot reload tools
- Optimized for performance and security

## Flow (Where it comes from -> Where it goes)

Dockerfile + env profiles -> docker compose -> running services

## Why It Was Designed This Way

- Keep runtime profiles isolated (dev vs prod).
- Keep Docker concerns out of application code.
- Use reproducible builds and consistent startup.
- Hot reload in dev, optimized binaries in prod.

## Recommended Practices Visible Here

- Multi-stage build for small images.
- Per-environment env files and compose overlays.
- Entry-point scripts are simple and explicit.
- Separate Dockerfiles for dev (with tools) vs prod (optimized).

## Common Commands

```bash
make dev           # Start dev environment (hot reload enabled)
make rebuild-api   # Rebuild API container
make logs-api      # View API logs (watch Air rebuilds)
```

## What Should NOT Live Here

- Business logic or app behavior.
- Real secrets committed to git.
- Environment-specific overrides in shared templates.
