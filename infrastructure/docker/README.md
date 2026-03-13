# Docker Infrastructure

**Path:** `infrastructure/docker`

## Overview

Docker build/runtime assets for local and production-like environments.
This package defines container images, compose profiles, and environment wiring.

## What It Contains

- root `Dockerfile`: production-oriented API image build (multi-stage).
- `environments/`: profile-specific compose/env definitions (`dev`, `prod`, `example`).
- `scripts/`: utility helpers for build/run/clean flows.

## Main Flows

- Dev: `make dev` (or `make dev-up` / `make clean`) using `environments/dev/docker-compose-dev.yaml` and `.env.dev`.
- Prod-like: `make prod` using `environments/prod/docker-compose-prod.yaml` and `.env.prod`.
- Custom profile: start from `environments/example/.env.example`.

## Design Notes

- Keep environment-specific details inside profile folders.
- Keep production and development image concerns separated.
- Keep compose/runtime behavior reproducible via Make targets.

## Dev Hot Reload Notes

The dev compose profile is wired for hot reload across integrated services:

- API: `infrastructure/docker/environments/dev/Dockerfile.dev` (Air)
- Chat: source mount + Uvicorn reload in dev mode
- Dashboard: source mounts + Vite HMR
- Kafka backbone via Redpanda-compatible local broker
- `aion-ingest` and `aion-streams` bootstrap services in the same local network

Operational commands:

```bash
make dev
make dev-fast
make rebuild-api
make rebuild-chat
make rebuild-dashboard
```

Use targeted rebuild commands when dependency layers or Dockerfiles changed.

## Best Practices

- Do not store real secrets in `.env.*`; keep secrets outside git.
- Keep service parity between dev and prod; change only tuning/security details.
- Use `make build-dev` / `make build-prod` for reproducible images (BuildKit enabled).
- Preserve healthchecks (Postgres/Redis/API) when adding services.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
