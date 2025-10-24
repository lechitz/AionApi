# AionApi — Documentation

Welcome to the AionApi documentation — a concise reference and walkthrough for this modular Go backend that powers habit and diary management with REST and GraphQL APIs.

---

## Build better habits with data

Aion helps you capture daily actions, measure patterns, and foster sustainable routines using a simple, extensible API-first platform. The server is designed for reliability, observability, and developer ergonomics — ideal for experimentation and iteration.

Key ideas:
- Instrumented by default (OpenTelemetry + Prometheus) so you can observe behavior and performance
- Hexagonal (Ports & Adapters) architecture that keeps business logic decoupled from transport and infrastructure
- Developer-first workflows with codegen, generated mocks, formatters and linters

---

## Quick links
- Getting started — Quick setup and local developer workflow: `docs/getting-started.md`
- Architecture — Design, component responsibilities and request flows: `docs/architecture.md`
- Platform — Observability, server and configuration notes: `docs/platform.md`
- API (Swagger) — interactive API explorer: https://lechitz.github.io/AionApi/swagger-ui/
- Raw OpenAPI spec: `./swagger/swagger.yaml`

---

## Quick start (dev)

These commands get you from repo clone to a running dev environment.

```bash
# clone
git clone git@github.com:lechitz/AionApi.git
cd AionApi

# install recommended developer tools (goimports, golines, golangci-lint, migrate, gqlgen, mockgen, ...)
make install-tools

# fetch modules
go mod download

# start the local development stack (Postgres, API, etc.)
make dev

# apply migrations and optionally seed sample data
export MIGRATION_DB="postgres://aion:aion@localhost:5432/aionapi?sslmode=disable"
make migrate-up
make seed-all

# health check
curl -s http://localhost:8080/aion/health | jq
```

> Note: `make install-tools` is a convenient alias for the repository tooling installer (`tools-install` in `makefiles/tooling.mk`).

---

## Explore the API
- Live interactive API (Swagger UI): https://lechitz.github.io/AionApi/swagger-ui/
- OpenAPI spec (raw): `https://raw.githubusercontent.com/lechitz/AionApi/main/swagger/swagger.yaml`

---

## Preview this documentation locally

If you want to preview the MkDocs site while editing:

```bash
# ensure mkdocs + mkdocs-material are installed in your environment
mkdocs serve
# open http://127.0.0.1:8000 in the browser
```

---
