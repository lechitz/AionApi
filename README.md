# AionApi

A modular Go backend exposing REST and GraphQL APIs with Hexagonal/Clean Architecture principles.

## Quick Links

- Docs site: [AionApi GitHub Pages](https://lechitz.github.io/AionApi/)
- Swagger UI: [AionApi Swagger](https://lechitz.github.io/AionApi/swagger-ui/)

## Project Overview

AionApi provides backend capabilities for habit/diary management workflows, including user management, auth, categories, tags, records, and chat integrations.

## Core Stack

- Go
- Chi (HTTP routing)
- gqlgen (GraphQL)
- PostgreSQL + GORM
- Redis
- OpenTelemetry + Grafana + Prometheus + Loki
- Docker / Docker Compose

## Common Workflows

```bash
make dev
make migrate-up
make seed-all
make test
make verify
```

<!-- docs-index:start -->
## Documentation Index

Repository README map by area.

### cmd
- [`cmd/README.md`](./cmd/README.md)
- [`cmd/api/README.md`](./cmd/api/README.md)

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

- Add architecture decision record links for major platform/domain choices.
- Add contributor quick-start for first local run and troubleshooting.
- Add release notes section summarizing contract changes per version.
- Add badges/automation links for CI, coverage, and docs publishing.
