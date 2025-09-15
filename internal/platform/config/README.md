# Config (Platform)

**Folder:** `internal/platform/config`

## Responsibility

* Provide a **typed, validated** configuration for the whole app (platform + domain).
* **Load** values from environment variables (via `envconfig`), **normalize** them, and **fail fast** on invalid input.
* **Bridge secrets** responsibly: generate a **temporary secret** if none is provided (for local/dev only).

---

## How it works

* `Loader` reads env vars with [`envconfig`](https://github.com/kelseyhightower/envconfig), builds a `Config`, validates it, and applies small normalizations.
* If `SECRET_KEY` is **missing**, the loader asks an injected **key generator** to produce an ephemeral key and logs a warning (never for prod).
* Timeouts are clamped by constants in `0_platform_config_constants.go`, e.g.:

    * `MinHTTPTimeout`
    * `MinGraphQLTimeout`

### Construction & usage

```go
import (
  "github.com/lechitz/AionApi/internal/platform/config"
  "github.com/lechitz/AionApi/internal/platform/ports/output/keygen"
  "github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

func Bootstrap(log logger.ContextLogger, kg keygen.Generator) (config.Config, error) {
  ld := config.NewLoader(log, kg) // constructor in loader.go
  cfg, err := ld.Load()
  if err != nil {
    return config.Config{}, err
  }
  return cfg, nil
}
```

---

## What it defines

`config.Config` aggregates all sections the platform exposes (see `config.go` / `environments.go`). Typical sections include:

* **General**

    * `Name`, `Env`, `Version`
    * Tags: `envconfig:"APP_NAME"`, `envconfig:"APP_ENV"`, `envconfig:"APP_VERSION"`
* **Secret**

    * `Key` (reads `SECRET_KEY`); generated if empty (dev convenience).
* **Observability**

    * OTEL exporter endpoint, service name/version, sampling settings, etc. (e.g., `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_SERVICE_NAME`).
* **HTTP Server**

    * Host/port, read/write timeouts, graceful `ShutdownTimeout`.
    * `Context` (base path) is **normalized** to start with `/` and to avoid trailing `/`.
* **GraphQL Server**

    * Path, timeouts, and toggles like `PlaygroundEnabled` (when present).
* **Database**

    * Either a DSN or discrete settings (host, port, user, pass, db, sslmode, pool limits).
* **(Other platform sections as needed)**

    * Cache, JWT, cookies, CORS, etc.—each with explicit `envconfig` tags.

> Exact field names/types live in `config.go` and `environments.go`. This readme documents the intent and conventions so new sections follow the same pattern.

---

## Validation & normalization

* **Timeouts**: values below `MinHTTPTimeout` / `MinGraphQLTimeout` are rejected.
* **HTTP/GraphQL context path**: normalized to `"/<prefix>"` (no trailing slash).
* **Required fields**: DB connectivity and other critical settings must be present (or resolvable from defaults). Missing required values return an error.
* **Secrets**:

    * If `SECRET_KEY` is missing, a **temporary** key is generated with the injected `keygen.Generator`.
    * A **warning** is logged; this is acceptable for local/dev, but you **must** set a real key for prod.

---

## Environment variables (examples)

```env
# General
APP_NAME=AionApi
APP_ENV=development
APP_VERSION=0.1.0

# Observability
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
OTEL_SERVICE_NAME=AionApi
OTEL_SERVICE_VERSION=0.1.0

# HTTP
SERVER_HTTP_CONTEXT=/aion
SERVER_HTTP_PORT=8080
SHUTDOWN_TIMEOUT=5

# GraphQL (examples)
GRAPHQL_PATH=/graphql
GRAPHQL_TIMEOUT=5s
GRAPHQL_PLAYGROUND=true

# Database (examples)
DB_HOST=localhost
DB_PORT=5432
DB_USER=aion
DB_PASSWORD=secret
DB_NAME=aionapi
DB_SSLMODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_CONN_MAX_LIFETIME=30m

# Secrets
SECRET_KEY=super-long-random-hex
```

> The exact var names match the `envconfig` tags in the structs; use the above as guidance.

---

## Design notes

* **Single source of truth**: only the config package should know env var names and normalization rules.
* **No side effects** beyond reading env vars, generating an ephemeral dev secret, and logging validation issues.
* Keep **platform concerns** (HTTP/GraphQL/DB/OTEL) here; domain packages receive values already parsed & validated.

---

## Testing hints

* Use `t.Setenv("VAR", "value")` to isolate each test.
* Inject a **fake keygen** for deterministic behavior when `SECRET_KEY` is absent.

  ```go
  kg := mocks.NewMockGenerator(ctrl)
  kg.EXPECT().Generate().Return("insecure-dev-key", nil)
  cfg, err := config.NewLoader(log, kg).Load()
  ```
* Assert:

    * normalization of `SERVER_HTTP_CONTEXT`
    * min timeout validation
    * error on incomplete DB settings (if required by your build)
    * warning path when secret is generated

---

## Gotchas

* Don’t rely on **implicit defaults** in production—set env vars explicitly.
* If you change migration paths, ports, or OTEL endpoints in Docker, reflect them in your `.env.*` files so `Loader` picks them up.
* Keep new fields **tagged with `envconfig`** and add validation if they’re critical.
