# AionApi Documentation

AionApi is a modular Go backend for habit and diary workflows, built with Hexagonal Architecture and strong observability by default.

## What You Can Do Here

| Area | What you will find |
| --- | --- |
| Start Here | Local setup, first run, and basic validation |
| Architecture | Layer boundaries, context isolation, and request lifecycle |
| API | REST contract (Swagger), GraphQL usage, and collections |
| Platform & Ops | Observability setup and operational guides |
| Reference | Diagram catalog, docs assets, and changelog |

## Quick Links

- Docs site: <https://lechitz.github.io/AionApi/>
- Swagger UI: <https://lechitz.github.io/AionApi/swagger-ui/>
- OpenAPI spec: <https://raw.githubusercontent.com/lechitz/AionApi/main/swagger/swagger.yaml>
- Repository: <https://github.com/lechitz/AionApi>

## Recommended Reading Path

1. Read [Getting Started](getting-started.md) and boot the local stack.
2. Read [System Design](architecture.md) to understand boundaries and flow.
3. Read [Platform Runtime](platform.md) for wiring, server, and observability internals.
4. Use [GraphQL Guide](graphql-api.md) and Swagger UI for API integration.

## Documentation Principles

- Architecture-first: docs follow the same boundaries as code.
- Operationally useful: commands are runnable in a fresh local environment.
- Contract-driven: API docs point to generated artifacts and source contracts.
- Maintainable: each page has a clear ownership scope and avoids duplicate low-level details.

## Keep This Portal Updated

When you change behavior, update docs in the same PR for:

- API contract changes (REST/GraphQL)
- Architecture changes (ports, adapters, context boundaries)
- Operational changes (make targets, observability stack, environment variables)

!!! tip
    This site should stay high-level and navigable. Deep package-level details remain in repository `README.md` files and code comments.
