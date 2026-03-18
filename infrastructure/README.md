# Infrastructure Layer

**Path:** `infrastructure`

## Purpose

This layer owns repo-local infrastructure assets used to run, provision, and observe `AionApi`.
It supports the multi-repo local stack, prod-like compose profiles, and database lifecycle work.

## Current Areas

| Area | Responsibility |
| --- | --- |
| `assets/` | static support assets published through the local asset bucket |
| `db/` | schema migrations and local seed datasets |
| `docker/` | Docker image build, compose profiles, and entrypoint wiring |
| `observability/` | OTel, Prometheus, Loki, Fluent Bit, and Grafana configs |

## Boundaries

- keep business logic out of `infrastructure`
- treat versioned SQL, compose files, and telemetry configs as code
- cross-repo orchestration still depends on sibling repos in the workspace; this folder owns only the `AionApi` side of that wiring
- secrets and machine-specific values belong in ignored env files, not committed infra docs

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../README.md)
<!-- doc-nav:end -->
