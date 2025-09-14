# Database Seeds

**Folder:** `infrastructure/db/seed`

## Responsibility

* Provide sample **users** and **tag categories** for local/dev databases.
* Make it easy to demo the API and run tests against realistic data.
* Keep seed data **human-readable** and **versioned** alongside the codebase.

---

## How it works

* Raw SQL files are executed **against Postgres** using `psql`.
* In dev, we seed an already-running container (default **`postgres-dev`**).
* Seeds are designed for a **fresh database** (they are **not** idempotent).

> Migrations live in `infrastructure/db/migrations` and create the schema.
> Seeds live here and insert example rows on top of that schema.

---

## Files

* `user.sql`
  Inserts five demo users (password hashes are placeholders suitable for dev).

* `category.sql`
  Inserts a few **tag categories**:

    * user **1**: `Reading`, `Health`, `Development`
    * user **2**: `Finance`

> Heads-up: `tag_categories.name` is **UNIQUE** (global), so re-running this script will fail unless you wipe data first.

---

## Running the seeds

From the **repo root**, with the dev stack up:

```bash
# Start/refresh dev Postgres
make dev-up

# Seed users
make seed-users

# Seed tag categories
make seed-categories

# Or everything
make seed-all
```

These make targets run:

* `docker exec -i postgres-dev psql -U aion -d aionapi < infrastructure/db/seed/user.sql`
* `docker exec -i postgres-dev psql -U aion -d aionapi < infrastructure/db/seed/category.sql`

> If you get “relation already exists” or “duplicate key” errors, you likely seeded a non-empty DB. See **Idempotency & safety** below.

---

## Idempotency & safety

* The scripts use plain `INSERT`. They **will** error or duplicate if run twice.
* To reseed, reset the database (e.g., **drop the volume**):

  ```bash
  make dev-down && make clean-dev && make dev-up && make seed-all
  ```

---

## Conventions

* Keep seeds **small** and **meaningful** (no PII, no real secrets).
* Prefer **deterministic** values (fixed emails/usernames) to simplify testing.
* If you add new tables, create a new seed file and (optionally) a new make target.

---

## Troubleshooting

* **Container not found**: Verify the container name is `postgres-dev` (see `docker ps`).
* **Path mismatch**: The root Makefile must reference `infrastructure/db/seed`.
  If it points to `infrastructure/db/seeds`, update it or rename this folder for consistency.
* **Auth failure**: Confirm DB creds in your `infrastructure/docker/environments/dev/.env.dev`.
* **Migrations missing**: Run migrations first (they are applied automatically on a fresh dev DB via the mounted `migrations` folder).
