# Category Repository (Secondary Adapter — DB)

**Path:** `internal/category/adapter/secondary/db/repository`

## Responsibilities

* Implement the **persistence output port** for the Category context (create, read, update, soft delete, list).
* **Hide ORM/SQL details** (GORM) behind domain types; expose only `domain.Category`.
* Provide **observability** (OTel spans, structured logs) for DB operations.

## How it works

* **Construction:** `NewCategory(db *gorm.DB, logger logger.ContextLogger) *CategoryRepository`
* **Domain ↔ DB mapping:** Uses `mapper.CategoryToDB/CategoryFromDB` and `model.CategoryDB`. Domain code never sees GORM types.
* **Context propagation:** All methods accept `context.Context` and call `db.WithContext(ctx)` to support cancellation, timeouts, and tracing.
* **Methods & semantics:**

    * `Create(ctx, category domain.Category) (domain.Category, error)`
      Inserts a new category; returns the created domain entity.
    * `GetByID(ctx, categoryID, userID uint64) (domain.Category, error)`
      Fetches by composite key; returns domain entity or an error on not found/DB failure.
    * `GetByName(ctx, name string, userID uint64) (domain.Category, error)`
      Fetches by `(user_id, name)`; returns **zero-value** `domain.Category` and `nil` when not found.
    * `ListAll(ctx, userID uint64) ([]domain.Category, error)`
      Returns all categories for a user.
    * `UpdateCategory(ctx, categoryID, userID uint64, updateFields map[string]any) (domain.Category, error)`
      Partial update; reads back the row and returns the updated domain entity.
    * `SoftDelete(ctx, categoryID, userID uint64) error`
      Sets `deleted_at`/`updated_at` (UTC) to mark as deleted.
* **Observability:** Each method opens an OTel span (`Tracer: "CategoryRepository"`), sets canonical attributes (e.g., `user_id`, `category_id`, `category_name`, `operation`), and records errors when present. Minimal, structured logging is used where helpful.

## My reminders

* **Do not leak** GORM models or DB errors to the domain/handlers; always return `domain.Category` and meaningful errors.
* **Consistency on “not found”:** pick one approach (zero-value + nil **or** typed error) and document it; services should normalize behavior.
* **Sanitize updates:** only allow intended fields in `updateFields`; never set immutable columns (e.g., `created_at`).
* **Use context everywhere:** pass `ctx`, attach attributes, and record errors on the span.
* **Log metadata only:** IDs and counts; never sensitive payloads.
* **Transactions:** make them explicit when a use case requires multiple statements with atomicity.
* **UTC timestamps** for soft deletes and updates.
* **Tests:** table tests for queries/edge cases; consider lightweight DB (e.g., SQLite) or containers for integration.
