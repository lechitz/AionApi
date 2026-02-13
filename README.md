# AionApi

AionApi is a production-oriented Go backend that exposes REST and GraphQL APIs for habit and diary workflows, built with Hexagonal/Clean Architecture and strong observability.

## Why This Project

AionApi focuses on three goals:

- keep business logic isolated from transport and infrastructure
- provide stable API contracts for multiple clients
- keep operations visible and debuggable in local and production-like environments

## Quick Links

- Documentation portal: [AionApi Docs](https://lechitz.github.io/AionApi/)
- REST explorer: [Swagger UI](https://lechitz.github.io/AionApi/swagger-ui/)
- OpenAPI contract: [`swagger/swagger.yaml`](./swagger/swagger.yaml)
- GraphQL schema artifact: [`docs/graphql/schema.graphql`](./docs/graphql/schema.graphql)

## Architecture At A Glance

| Layer | Purpose |
| --- | --- |
| `internal/<ctx>/core` | Domain, ports, and usecases (business logic) |
| `internal/<ctx>/adapter/primary` | HTTP/GraphQL input adapters |
| `internal/<ctx>/adapter/secondary` | DB/cache/provider output adapters |
| `internal/platform` | App bootstrap, server, config, and observability wiring |
| `infrastructure` | Docker, migrations, observability stack |

## Core Stack

- Go
- Chi (HTTP routing)
- gqlgen (GraphQL)
- PostgreSQL + GORM
- Redis
- OpenTelemetry + Prometheus + Grafana + Loki
- Docker / Docker Compose

## Fast Local Workflow

```bash
make tools-install
make dev
make migrate-up
make seed-all
make verify
```

## Quality Gates

```bash
make test
make test-cover-detail
make docs-verify
make verify
```

<!-- docs-index:start -->
## Documentation Index

Repository README map by area.

### cmd
- [`cmd/README.md`](./cmd/README.md)

### contracts
- [`contracts/graphql/queries/README.md`](./contracts/graphql/queries/README.md)
- [`contracts/openapi/README.md`](./contracts/openapi/README.md)

### docs
- [`docs/assets/README.md`](./docs/assets/README.md)
- [`docs/collections/README.md`](./docs/collections/README.md)
- [`docs/diagram/README.md`](./docs/diagram/README.md)
- [`docs/graphql/README.md`](./docs/graphql/README.md)
- [`docs/swagger-ui/README.md`](./docs/swagger-ui/README.md)

### hack
- [`hack/README.md`](./hack/README.md)
- [`hack/dev/README.md`](./hack/dev/README.md)
- [`hack/tools/seed-caller/README.md`](./hack/tools/seed-caller/README.md)
- [`hack/tools/seed-helper/README.md`](./hack/tools/seed-helper/README.md)

### infrastructure
- [`infrastructure/README.md`](./infrastructure/README.md)
- [`infrastructure/db/README.md`](./infrastructure/db/README.md)
- [`infrastructure/db/migrations/README.md`](./infrastructure/db/migrations/README.md)
- [`infrastructure/db/seed/README.md`](./infrastructure/db/seed/README.md)
- [`infrastructure/docker/README.md`](./infrastructure/docker/README.md)
- [`infrastructure/docker/environments/README.md`](./infrastructure/docker/environments/README.md)
- [`infrastructure/docker/environments/example/README.md`](./infrastructure/docker/environments/example/README.md)
- [`infrastructure/observability/README.md`](./infrastructure/observability/README.md)
- [`infrastructure/observability/fluentbit/README.md`](./infrastructure/observability/fluentbit/README.md)
- [`infrastructure/observability/grafana/README.md`](./infrastructure/observability/grafana/README.md)
- [`infrastructure/observability/loki/README.md`](./infrastructure/observability/loki/README.md)
- [`infrastructure/observability/otel/README.md`](./infrastructure/observability/otel/README.md)
- [`infrastructure/observability/prometheus/README.md`](./infrastructure/observability/prometheus/README.md)

### internal
- [`internal/README.md`](./internal/README.md)
- [`internal/adapter/README.md`](./internal/adapter/README.md)
- [`internal/adapter/primary/README.md`](./internal/adapter/primary/README.md)
- [`internal/adapter/primary/graphql/README.md`](./internal/adapter/primary/graphql/README.md)
- [`internal/adapter/secondary/README.md`](./internal/adapter/secondary/README.md)
- [`internal/admin/README.md`](./internal/admin/README.md)
- [`internal/auth/README.md`](./internal/auth/README.md)
- [`internal/category/README.md`](./internal/category/README.md)
- [`internal/chat/README.md`](./internal/chat/README.md)
- [`internal/platform/README.md`](./internal/platform/README.md)
- [`internal/platform/config/README.md`](./internal/platform/config/README.md)
- [`internal/platform/fxapp/README.md`](./internal/platform/fxapp/README.md)
- [`internal/platform/httpclient/README.md`](./internal/platform/httpclient/README.md)
- [`internal/platform/observability/README.md`](./internal/platform/observability/README.md)
- [`internal/platform/ports/README.md`](./internal/platform/ports/README.md)
- [`internal/platform/server/README.md`](./internal/platform/server/README.md)
- [`internal/platform/server/http/README.md`](./internal/platform/server/http/README.md)
- [`internal/platform/server/http/generic/README.md`](./internal/platform/server/http/generic/README.md)
- [`internal/platform/server/http/middleware/README.md`](./internal/platform/server/http/middleware/README.md)
- [`internal/platform/server/http/middleware/cors/README.md`](./internal/platform/server/http/middleware/cors/README.md)
- [`internal/platform/server/http/middleware/recovery/README.md`](./internal/platform/server/http/middleware/recovery/README.md)
- [`internal/platform/server/http/middleware/requestid/README.md`](./internal/platform/server/http/middleware/requestid/README.md)
- [`internal/platform/server/http/middleware/servicetoken/README.md`](./internal/platform/server/http/middleware/servicetoken/README.md)
- [`internal/platform/server/http/ports/README.md`](./internal/platform/server/http/ports/README.md)
- [`internal/platform/server/http/router/README.md`](./internal/platform/server/http/router/README.md)
- [`internal/platform/server/http/utils/README.md`](./internal/platform/server/http/utils/README.md)
- [`internal/platform/server/http/utils/cookies/README.md`](./internal/platform/server/http/utils/cookies/README.md)
- [`internal/platform/server/http/utils/httpresponse/README.md`](./internal/platform/server/http/utils/httpresponse/README.md)
- [`internal/platform/server/http/utils/sharederrors/README.md`](./internal/platform/server/http/utils/sharederrors/README.md)
- [`internal/record/README.md`](./internal/record/README.md)
- [`internal/shared/README.md`](./internal/shared/README.md)
- [`internal/shared/constants/README.md`](./internal/shared/constants/README.md)
- [`internal/tag/README.md`](./internal/tag/README.md)
- [`internal/user/README.md`](./internal/user/README.md)

### makefiles
- [`makefiles/README.md`](./makefiles/README.md)

### tests
- [`tests/coverage/README.md`](./tests/coverage/README.md)
- [`tests/setup/README.md`](./tests/setup/README.md)

<!-- docs-index:end -->

## Package Improvements

- Add architecture decision records (ADRs) for critical platform/domain choices.
- Add release notes summary per version with API contract deltas.
- Add contributor troubleshooting matrix for local setup failures.
- Add CI/docs badges and links to pipeline checks.
