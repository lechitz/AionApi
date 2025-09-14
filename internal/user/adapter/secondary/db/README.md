# User DB Adapter (Secondary)

**Folder:** `internal/user/adapter/secondary/db`

## Responsibility

* Implement the **persistence output port** for the **User** context.
* Hide **ORM/SQL** (GORM) behind domain types; expose only `domain.User`.
* Provide **observability** (OpenTelemetry spans + structured logs) for every DB call.
* Enforce **soft-delete semantics** and **uniqueness** rules at the repository boundary.

## How it works

* **Construction & DI:** The repository is created during platform bootstrap and injected into the user use case via the **output port**. It holds a `*gorm.DB` and a `logger.ContextLogger`.
* **Tables:** Uses `aion_api.users` (see migrations). Soft-deleted rows have `deleted_at` set; reads must ignore them.
* **Context-first:** Every method accepts a `context.Context` (cancellation, deadlines, trace propagation).
* **No business logic here:** Hashing, validation, and policies live in the **use case**; the repo only persists and maps data.

## Operations (output port)

* `CheckUniqueness(ctx, username, email)` — probes DB for username/email conflicts and returns a `UserUniqueness` view.
* `Create(ctx, user)` — inserts a new user and returns the persisted `domain.User`.
* `GetByID(ctx, userID)` — fetches a user by ID (ignores soft-deleted).
* `GetByUsername(ctx, username)` — fetches by username (ignores soft-deleted).
* `GetByEmail(ctx, email)` — fetches by email (ignores soft-deleted).
* `ListAll(ctx)` — returns all non-deleted users (consider pagination in callers).
* `Update(ctx, userID, fields)` — partial update (only provided fields).
* `SoftDelete(ctx, userID)` — marks the user as deleted by setting `deleted_at`.

> All reads must add `WHERE deleted_at IS NULL`. Updates/soft-deletes should check **rows affected** and return a not-found semantic error when zero.

## Mapping rules

* **DB ⇄ Domain:** Implemented in `mapper/user_mapper.go`:

    * `model.UserDB` → `domain.User` (IDs, strings, and timestamps).
    * Respect `gorm.Model` timestamps when present; ensure `DeletedAt` is omitted from domain unless needed by a read model.
* **Never** return DB structs outside the adapter. Domain types only.

## Observability

* Tracer & span names in `repository/0_user_repository_constants.go` (e.g., `TracerUserRepository`, `SpanCreate`, `SpanGetByID`, …).
* Add canonical attributes such as:

    * `user_id`, `username`, `email` (when available)
    * `db.table="aion_api.users"`, `db.operation` (`select|insert|update`)
* Record errors with `span.RecordError` and set `codes.Error` on failure; set `codes.Ok` on success.

## Error semantics

* Translate driver/DB errors to **semantic** errors used by the domain:

    * **Not found** → return a domain “not found” error.
    * **Uniqueness** (username/email) → return a domain **validation**/conflict error.
    * **Precondition/empty update** → validation error.
* Do **not** leak raw SQL/driver messages to callers or logs. Log **metadata**, never secrets.

## Transactions

* Single-row operations run in the ambient connection.
* If a use case needs multi-step atomicity, it should orchestrate a transaction and pass the transactional context/DB handle to the repository (future extension).

## Files

* **`model/user_model.go`** — GORM model + `TableUsers` constant.
* **`mapper/user_mapper.go`** — DB ⇄ domain mapping helpers.
* **`repository/0_user_repository_impl.go`** — repository struct (DB + logger) and shared wiring.
* **`repository/0_user_repository_constants.go`** — tracer/span names.
* **`repository/*.go`** — per-operation implementations:

    * `create.go`, `update.go`, `soft_delete.go`
    * `get_by_id.go`, `get_by_username.go`, `get_by_email.go`, `list_all.go`
    * `check_uniqueness.go`

## Conventions & gotchas

* **Never hash passwords here.** The use case must pass hashed values.
* Always **filter soft-deleted** rows on reads.
* Keep **partial updates** safe: only update fields present in `fields`.
* Be strict about **case sensitivity** for username/email if the product requires it (consider `LOWER()` normalization at write/read if needed).

## Testing hints

* Use a **mocked logger** and an in-memory/isolated DB (or test transaction).
* Assert:

    * correct **span** creation and attributes
    * **uniqueness** behavior for username/email
    * **soft-delete** exclusion on reads
    * **rows affected** handling on update/soft-delete
* Seed with deterministic fixtures when validating list/order expectations.

## Future work

* Add **pagination & filtering** to `ListAll`.
* Introduce **scoped queries** (e.g., by tenant/org) if/when multi-tenancy appears.
* Consider **repository interface** consolidation for easier mocking where needed.
