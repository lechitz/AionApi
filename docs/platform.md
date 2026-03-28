# Platform Runtime

This page documents the cross-cutting runtime layer: configuration, bootstrap wiring, HTTP/GraphQL server composition, and observability setup.

## Platform Surface

```text
internal/platform/
  bootstrap/
  config/
  server/
    http/
    graph/
  observability/
```

## Configuration Model

Configuration is loaded from environment variables into a typed struct with validation and normalization.

| Area | Examples |
| --- | --- |
| HTTP | `SERVER_HTTP_PORT`, `SERVER_HTTP_CONTEXT`, `SHUTDOWN_TIMEOUT` |
| GraphQL | `GRAPHQL_PATH`, `GRAPHQL_TIMEOUT`, `GRAPHQL_PLAYGROUND` |
| Database | `DB_HOST`, `DB_PORT`, `DB_MAX_OPEN_CONNS` |
| Cache | `CACHE_ADDR`, `CACHE_DB` |
| Observability | `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_SERVICE_NAME` |

## Bootstrap (Composition Root)

Bootstrap builds concrete adapters and returns application dependencies plus cleanup logic.

Responsibilities:

1. Connect Postgres/Redis and initialize providers.
2. Instantiate repositories/stores implementing output ports.
3. Build usecases implementing input ports.
4. Return `AppDependencies` and graceful `cleanup(ctx)`.

## HTTP Runtime

HTTP uses a router port abstraction with a `chi` implementation.

| Part | Responsibility |
| --- | --- |
| Router port | Framework-agnostic contract |
| Middleware | Recovery, request ID, auth and policy chain |
| Composer | Route mounting and global defaults |
| Generic handlers | Health, 404/405, and centralized error mapping |

### Middleware Order

Keep recovery outermost and request-id early to preserve traceability across logs and spans.

## GraphQL Runtime

GraphQL is mounted in platform server and delegates to context controllers/resolvers.

- Schema modules live under `internal/adapter/primary/graphql/schema/modules/`.
- Regenerate runtime artifacts with `make graphql`.
- Keep controllers thin and usecases authoritative.

## Observability Runtime

| Signal | Stack |
| --- | --- |
| Traces | OpenTelemetry + OTLP exporter + Jaeger |
| Metrics | Prometheus + Grafana dashboards |
| Logs | Structured logs + Fluent Bit/Loki integration |

Recommended local environment:

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"
export OTEL_SERVICE_NAME="aion-api"
export OTEL_SERVICE_VERSION="0.0.1"
```

## Operational Commands

```bash
make dev-up
make dev-down
make verify
make docs-build
```

## Next Step

- API usage: [GraphQL Guide](graphql-api.md)
- Operational validation: [Observability Quickstart](observability-quickstart.md)
