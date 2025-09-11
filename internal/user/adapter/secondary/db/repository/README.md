# User Repository (Secondary Adapter — DB)

**Path:** `internal/user/adapter/secondary/db/repository`

## Responsibilities

* Implement the **persistence output port** of the **User** context.
* Hide **GORM/SQL** details behind domain types; expose only `domain.User`.
* Enforce **soft-delete** semantics and **uniqueness checks** (username/email).
* Emit **observability** signals (OpenTelemetry spans + structured logs).

## How it works

* **Construction & DI:** The repository is instantiated at bootstrap with:

    * `*gorm.DB` (connection/transaction handle)
    * `logger.ContextLogger`
* **Context-first API:** Every method receives a `context.Context` for cancellation, deadlines, and trace propagation.
* **No business rules here:** Hashing/validation/policies live in the **use case**. The repository persists and maps.

## Operations

* `CheckUniqueness(ctx, username, email)`
  Looks for username/email conflicts and returns a `UserUniqueness` view for the use case to decide.
* `Create(ctx, user)`
  Inserts a new user record and returns the persisted `domain.User`.
* `GetByID(ctx, userID)`
  Fetches an existing user by id (**ignores soft-deleted** rows).
* `GetByUsername(ctx, username)` / `GetByEmail(ctx, email)`
  Fetch by alternative keys (**ignores soft-deleted** rows).
* `ListAll(ctx)`
  Returns all non-deleted users (consider pagination upstream).
* `Update(ctx, userID, fields)`
  Partial update; only provided fields are persisted.
* `SoftDelete(ctx, userID)`
  Sets `deleted_at` (soft delete). Rows remain in DB for audit/recovery.

> Reads must add `WHERE deleted_at IS NULL`. Updates/soft-deletes should verify **rows affected** and surface a **not found** error when zero.

## Mapping

* Mapping between DB and domain lives in `internal/user/adapter/secondary/db/mapper`:

    * `model.UserDB` ⇄ `domain.User`
* Never expose GORM models outside this adapter.

## Observability

* Tracer/span names in `0_user_repository_constants.go` (e.g., `TracerUserRepository`, `SpanCreate`, `SpanGetByID`, …).
* Add canonical span attributes such as:

    * `user_id`, `username`, `email` (when present)
    * `db.table="aion_api.users"`, `db.operation` (`select|insert|update`)
* Record failures with `span.RecordError(err)` and set `codes.Error`; set `codes.Ok` on success.
* Log **metadata only** (no secrets, no raw SQL).

## Error semantics

Translate driver/GORM errors to **semantic** domain errors:

* **Not found** → domain not-found error.
* **Uniqueness violation** (username/email) → domain validation/conflict error (detected in `CheckUniqueness` and during inserts/updates).
* **Empty/invalid patch** → domain validation error.
* Do not leak driver-specific messages to the domain.

## Transactions

* Single-step operations use the ambient connection.
* Multi-step atomic flows should be coordinated by the **use case**, passing a transactional handle (future extension if needed).

## Files

* **`0_user_repository_impl.go`** — struct/wiring (`*gorm.DB`, `logger`), shared helpers.
* **`0_user_repository_constants.go`** — tracer/span names.
* **`create.go`**, **`update.go`**, **`soft_delete.go`** — write operations.
* **`get_by_id.go`**, **`get_by_username.go`**, **`get_by_email.go`**, **`list_all.go`** — read operations.
* **`check_uniqueness.go`** — username/email conflict probing.

## Conventions & gotchas

* **Never hash passwords here**; the use case must pass hashed values.
* Always exclude soft-deleted rows from reads.
* Keep partial updates **explicit**: update only keys present in `fields`.
* Be deliberate about username/email **case sensitivity**; normalize at boundaries if required by product rules.

## Testing hints

* Use a mock logger and a test DB/transaction for isolation.
* Assert for each method:

    * Correct span creation + attributes
    * Uniqueness detection on username/email
    * Soft-delete exclusion on reads
    * Proper “rows affected” handling for update/soft-delete
* Seed deterministic fixtures when verifying list ordering.

## Future work

* Add pagination & filtering to `ListAll`.
* Surface repository interface for easier mocking (if you want to mock at the persistence boundary).
* Optional read models that include `deleted_at` for admin/audit use cases.
