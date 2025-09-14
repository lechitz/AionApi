# Docker

**Folder:** `infrastructure/docker`

## Responsibility

* Containerize the **Aion API** and its infra dependencies.
* Separate **dev** and **prod** orchestration with clear, environment-specific config.
* Provide a repeatable workflow for **building images**, **running compose**, and **bootstrapping the DB** (migrations).

---

## Layout

```
infrastructure/docker/
├── Dockerfile                     # Multi-stage build (Go -> Alpine)
├── environments/
│   ├── dev/
│   │   ├── docker-compose-dev.yaml
│   │   └── .env.dev               # NOT for production; developer defaults
│   ├── example/
│   │   └── .env.example           # Template of required vars
│   └── prod/
│       ├── docker-compose-prod.yaml
│       └── .env.prod              # Example of production variables (do not commit secrets)
└── scripts/
    └── entrypoint.sh              # Starts the compiled binary in the runtime image
```

---

## How it works

* **Dockerfile (multi-stage):**

    * Stage 1: builds a static `aion-api` binary from `./cmd/aion-api`.
    * Stage 2: minimal **Alpine** runtime that runs `./aion-api` via `scripts/entrypoint.sh`.

* **Compose (dev/prod):**

    * Brings up **Postgres** with a persistent volume.
    * On first boot, Postgres executes any `.sql` files mounted at `/docker-entrypoint-initdb.d`.
    * We mount the project’s migrations directory so the schema is initialized automatically:

        * Dev: `infrastructure/db/migrations` → `/docker-entrypoint-initdb.d`

* **Environment files (`.env.*`):**

    * Consumed by `docker-compose` (`env_file:`) and Make targets.
    * Keep **secrets out of VCS**. Use `.env.example` as a reference and provide real values locally/CI.

> Today the provided compose files focus on **Postgres**. The API container is built with the Dockerfile and can be added to compose when desired (see “Extending compose with the API” below).

---

## Common workflows (via root Makefile)

* **Build DEV image**

  ```bash
  make build-dev
  ```

* **Start DEV stack (compose)**

  ```bash
  # Exports variables from environments/dev/.env.dev and runs docker compose
  make dev-up
  ```

* **Stop & remove DEV stack (volumes included)**

  ```bash
  make dev-down
  ```

* **Clean leftover DEV resources (containers/volumes/images)**

  ```bash
  make clean-dev
  ```

* **Build/Start PROD stack**

  ```bash
  make build-prod
  make prod-up
  ```

* **Global clean**

  ```bash
  make docker-clean-all
  ```

---

## Dev specifics

* **Postgres**

    * Container name: `postgres-dev`
    * Port: `5432:5432` (localhost exposure)
    * Persistent volume: `postgres-data-dev`
    * **Migrations:** `infrastructure/db/migrations` is mounted into `/docker-entrypoint-initdb.d` so tables, FKs, triggers, etc. are created at first boot.

* **Typical local DSN**

  ```
  postgres://aion:<password>@localhost:5432/aionapi?sslmode=disable
  ```

---

## Prod specifics

* **.env.prod** holds the **shape** of production settings (observability, server, etc.).
  Do not store real secrets in Git. Use your platform’s secret store or CI/CD variable injection.

* **Volumes and ports** are defined in `docker-compose-prod.yaml`. Adjust mounts so the DB can be initialized with the same `infrastructure/db/migrations` folder (path must be correct relative to that file).

---

## Extending compose with the API

If/when you want the API to run as a container alongside Postgres, add a service similar to:

```yaml
services:
  aion-api:
    image: aion-api:dev            # or :prod after make build-*
    container_name: aion-api-dev
    depends_on:
      - postgres
    env_file:
      - ./.env.dev                 # or .env.prod
    command: ["/bin/sh", "-c", "/scripts/entrypoint.sh"]
    volumes:
      - ../../scripts:/scripts
    ports:
      - "8080:8080"                # if your HTTP server runs on 8080
    restart: on-failure
```

> The multi-stage image already contains the compiled `aion-api` binary. The `entrypoint.sh` simply executes it.

---

## Tips & troubleshooting

* **Fresh DB init:** remove the Postgres volume so `/docker-entrypoint-initdb.d` scripts run again:

  ```bash
  make dev-down && make clean-dev && make dev-up
  ```
* **Migrations didn’t run?** Check the **relative mount path** from your composition file to `infrastructure/db/migrations`.
* **Port in use?** Ensure nothing else is bound to `5432` (or change the host port).
* **Logs:** `docker compose logs -f` from the environment folder is your friend.

---

## Conventions

* Keep images small (multi-stage, static binary).
* Keep env files **scoped per environment**; never commit real secrets.
* Prefer **infrastructure/db/migrations** as the single source of truth for schema bootstrap.
* Compose files should remain thin and environment-specific; use the root Make targets to orchestrate.
