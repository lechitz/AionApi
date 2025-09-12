# Platform

This page describes the cross‑cutting platform layer that powers AionApi: configuration, dependency bootstrap, HTTP/GraphQL servers, middlewares, observability, and Docker environments. If you’re wiring a new context (feature) or deploying locally, start here.

> TL;DR
>
> * Typed **Config** from env vars, with sane defaults and validation
> * Central **Bootstrap** that wires DB, cache, token, hasher, logger
> * Framework‑agnostic **Router Port** + **chi adapter**
> * First‑class **Observability** (OpenTelemetry traces/metrics, Prometheus, Grafana)
> * Clear separation of concerns: **platform** vs. **domain**

---

## 1) Directory map (platform layer)

```
internal/platform/
  bootstrap/         # build infra adapters + services, return AppDependencies + cleanup
  config/            # env → typed config + validation/normalization
  server/
    http/            # router port, chi adapter, generic handlers, middlewares, composer
    graph/           # GraphQL wiring (schema assembly via make graphql)
  observability/     # OTel tracer/meter setup (OTLP/HTTP), helpers
```

Complementary infra lives under `infrastructure/` (Docker, migrations, OTEL, Prometheus, Grafana, Loki/Fluent Bit).

---

## 2) Configuration

The `config` package loads environment variables into a typed struct, validates values, and normalizes things like the HTTP context path.

### 2.1 Quick example (env)

```env
# General
APP_NAME=AionApi
APP_ENV=development
APP_VERSION=0.1.0

# HTTP
SERVER_HTTP_CONTEXT=/aion-api
SERVER_HTTP_PORT=8080
SHUTDOWN_TIMEOUT=5s

# GraphQL
GRAPHQL_PATH=/graphql
GRAPHQL_TIMEOUT=5s
GRAPHQL_PLAYGROUND=true

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=aion
DB_PASSWORD=secret
DB_NAME=aionapi
DB_SSLMODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_CONN_MAX_LIFETIME=30m

# Cache (Redis)
CACHE_ADDR=localhost:6379
CACHE_DB=0
CACHE_PASSWORD=

# Secrets
SECRET_KEY=super-long-random-hex

# Observability (OTLP over HTTP)
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
OTEL_SERVICE_NAME=AionApi
OTEL_SERVICE_VERSION=0.1.0
OTEL_EXPORTER_OTLP_HEADERS=x-api-key=dev
OTEL_INSECURE=true
```

!!! tip
In **dev**, if `SECRET_KEY` is missing, the loader may generate a **temporary key** and log a warning. Always set a real key in production.

### 2.2 Guarantees

* Minimal **timeouts** are enforced (HTTP/GraphQL) to avoid foot‑guns.
* `SERVER_HTTP_CONTEXT` is normalized to `"/prefix"` with no trailing slash.
* Config contains only **platform‑level** concerns; domain code receives values already parsed and validated.

---

## 3) Bootstrap (composition root)

The `bootstrap` package constructs infrastructure adapters and wires **usecases** (services) through domain **ports**.

**Responsibilities**:

1. Build **secondary adapters**: Postgres (GORM), Redis, token provider (e.g., JWT), password hasher (bcrypt/argon2), context logger (Zap), etc.
2. Instantiate repositories and stores.
3. Construct usecase services (Auth, User, Category, Tag, Admin) via **input ports**.
4. Return `*AppDependencies` and a `cleanup(ctx)` function to gracefully close DB/cache.

**Pseudocode**:

```go
// InitializeDependencies(ctx, cfg, log) (*AppDependencies, cleanup, error)
// - Connect DB (GORM) + Redis
// - New token provider + hasher + logger
// - New repositories/stores
// - New usecase services
// - return deps + cleanup
```

!!! note
The bootstrap **does not** mount HTTP routes or start servers. It only builds dependencies.

---

## 4) HTTP Server

The HTTP layer is framework‑agnostic thanks to a **Router Port**. The default implementation uses **chi**.

```
internal/platform/server/http/
  ports/                  # Router interface (framework‑agnostic)
  router/chi/             # chi adapter implementing the port
  generic/                # health, error/404/405 handlers
  middleware/
    requestid/            # ensures X-Request-ID
    recovery/             # catches panics, returns uniform 500
  composer.go             # wires routes, middlewares, defaults
  server.go               # *http.Server startup/shutdown
```

### 4.1 Router Port (contract)

```go
type Middleware func(http.Handler) http.Handler

type Router interface {
  Use(mw ...Middleware)
  Group(prefix string, fn func(Router))
  GroupWith(mw Middleware, fn func(Router))
  Mount(prefix string, h http.Handler)
  GET/POST/PUT/DELETE(path string, h http.Handler)
  SetNotFound(h http.Handler)
  SetMethodNotAllowed(h http.Handler)
  SetError(func(http.ResponseWriter, *http.Request, error))
  ServeHTTP(http.ResponseWriter, *http.Request)
}
```

### 4.2 Composition (excerpt)

```go
r := chiadapter.New()

// Global middlewares (order matters)
r.Use(recovery.New(genericHandler)) // outermost guard
r.Use(requestid.New())              // ensure X-Request-ID

// Platform defaults
r.SetNotFound(http.HandlerFunc(genericHandler.NotFoundHandler))
r.SetMethodNotAllowed(http.HandlerFunc(genericHandler.MethodNotAllowedHandler))
r.SetError(genericHandler.ErrorHandler)

// Context routes (examples)
userhandler.RegisterHTTP(r, userH)
userhandler.RegisterHTTPProtected(r, userH, authMiddleware)

// GraphQL
r.Mount(cfg.ServerGraphql.Path, graphqlHTTPHandler)
```

!!! warning
**Middleware order matters**: keep **recovery** first so it can catch panics from everything else. `requestid` should run early to propagate the same ID through logs and traces.

### 4.3 Generic endpoints

* `GET /health` — returns service metadata.
* **404/405** — standardized JSON bodies.
* **Error/Recovery** — uniform 500 via centralized error handler.

### 4.4 Error & response shape

Adapters use `internal/shared/httpresponse` + `sharederrors` to produce a consistent envelope and status mapping:

| Domain error       | HTTP status |
| ------------------ | ----------- |
| validation         | 400         |
| unauthorized       | 401         |
| forbidden          | 403         |
| not\_found         | 404         |
| conflict           | 409         |
| rate\_limited      | 429         |
| internal (default) | 500         |

---

## 5) GraphQL Server

GraphQL lives alongside HTTP and is mounted under `GRAPHQL_PATH` (default `/graphql`). Schema is **modular** per context.

**Flow**:

1. Drop `*.graphqls` files under each context (`internal/<ctx>/adapter/primary/graphql/schema`).
2. Run `make graphql` to copy modules into the platform graph package and execute **gqlgen**.
3. Resolvers delegate to context **handlers**, which call **usecases**.

```text
internal/platform/server/graph/schema/_modules/   # assembled by make graphql
```

!!! tip
Keep **business rules** in usecases; GraphQL handlers are thin mappers with tracing.

---

## 6) Observability

AionApi is **tracing‑first**. Traces and metrics are exported via **OTLP/HTTP** to an OpenTelemetry Collector.

### 6.1 Tracing

* **Tracer names** by layer, e.g.:

    * `aionapi.user.handler.http`
    * `aionapi.user.usecase`
    * `aionapi.user.repository`
    * `aionapi.graphql.handler`
* Each operation starts a span and sets canonical attributes (operation, ids, status). Errors call `span.RecordError(err)` and set error status.

### 6.2 Metrics

* Prometheus scraping configured in `infrastructure/observability/prometheus/prometheus.yml`.
* Grafana provisioning + example dashboard under `infrastructure/observability/grafana/`.

### 6.3 Collector config (dev)

* Default endpoint: `OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318`
* Optional headers via CSV: `OTEL_EXPORTER_OTLP_HEADERS="x-api-key=dev"`
* Set `OTEL_INSECURE=true` for plain HTTP in local setups.

!!! note
Logs are structured (Zap). Fluent Bit + Loki are optional and can be added later for log aggregation.

---

## 7) Docker & environments

Docker Compose files live in `infrastructure/docker/` (dev and prod‑like). Use the `Makefile` targets for convenience.

**Common targets**:

```bash
make dev-up      # start dev stack (API + DB + dependencies)
make dev-down    # stop & remove dev stack
make dev         # build image + dev-up
make prod        # build and run prod-like stack (where applicable)
```

Environment files are kept under `infrastructure/docker/environments/`. Review and align them with the **Config** variables above.

---

## 8) Security (HTTP/Auth)

* `Authorization: Bearer <token>` is validated by the **auth middleware** (primary adapter). On success, it injects `user_id` into `context` for downstream handlers.
* Cookie helpers (`internal/platform/server/http/helpers/httpresponse`) centralize secure cookie settings when a browser flow is used.
* **Never log secrets**: passwords, raw tokens, or cookie contents.

---

## 9) Adding a new context (platform checklist)

1. **Ports & Usecase**: define input/output ports and implement usecase under `internal/<ctx>/core`.
2. **Adapters**:

    * Primary: HTTP and/or GraphQL handler + DTO/schema
    * Secondary: repository (DB), cache, external clients
3. **Bootstrap**: wire repositories + services in `internal/platform/bootstrap`.
4. **HTTP composer**: mount routes via `RegisterHTTP(...)` and protected routes via `RegisterHTTPProtected(...)`.
5. **GraphQL**: add `schema/<ctx>.graphqls`, run `make graphql`.
6. **Docs**: add `docs/<ctx>.md` and update `mkdocs.yml` navigation.

---

## 10) Troubleshooting

**No traces/metrics show up**

* Check `OTEL_EXPORTER_OTLP_ENDPOINT` and that the Collector is listening on `/v1/traces` `/v1/metrics` over **HTTP**.
* If using custom headers, confirm `OTEL_EXPORTER_OTLP_HEADERS` CSV is valid.

**Chi routes return 404**

* Ensure the global HTTP **context prefix** is correct (`SERVER_HTTP_CONTEXT`, e.g., `/aion-api`). All routes are mounted beneath it.

**JWT accepted but user not authorized**

* Verify the **auth middleware** is applied to the subtree via `GroupWith(middleware.Auth, ...)`.
* Confirm token cache (AuthStore) contains the same reference stored at login.

**Migrations don’t run**

* Export `MIGRATION_DB` DSN and `MIGRATION_PATH` environment variables used by the `make migrate-*` targets.

**GraphQL schema changes not reflected**

* Run `make graphql` and then `go mod tidy`. Rebuild the service.

---

## 11) Further reading

* `docs/architecture.md` — big picture of Hexagonal + request flows
* `docs/getting-started.md` — local setup, make targets, common commands
* Source readmes under `internal/platform/*` for deep dives into each package
