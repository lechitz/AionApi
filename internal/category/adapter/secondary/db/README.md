# Category DB (Secondary Adapter)

**Folder:** `internal/category/adapter/secondary/db`

## Responsibility

* Implement the **persistence adapter** for the Category context behind the output port.
* Encapsulate **GORM/SQL** details; expose only **domain** types to the core.
* Provide **observability** for DB calls (OpenTelemetry spans + structured logging).
* Keep **mapping** between DB models and domain entities centralized and consistent.

---

## How it works

**Package layout**

* `model/` — GORM models that mirror the database schema.

    * `CategoryTable = "aion_api.tag_categories"`.
    * Timestamps and soft‐delete fields (`created_at`, `updated_at`, `deleted_at`) match the migrations.
* `mapper/` — pure functions to convert **DB ⇄ domain**:

    * `CategoryFromDB(model.CategoryDB) domain.Category`
    * (and the inverse when needed)
* `repository/` — concrete implementation of the **CategoryRepository** using `*gorm.DB` + logger.

    * Constructor: `NewCategory(db *gorm.DB, logger logger.ContextLogger) *CategoryRepository`
    * Methods (output port):

        * `Create(ctx, category) (domain.Category, error)`
        * `GetByID(ctx, id, userID) (domain.Category, error)`
        * `GetByName(ctx, name, userID) (domain.Category, error)`
        * `ListAll(ctx, userID) ([]domain.Category, error)`
        * `Update(ctx, id, userID, fields) (domain.Category, error)`
        * `SoftDelete(ctx, id, userID) error`
    * All DB calls use `db.WithContext(ctx)` and return **domain** types.

**DB schema source of truth**

* Migrations live under `infrastructure/db/migrations`.
* Category rows are stored in `aion_api.tag_categories` (see `model.CategoryTable`).

---

## Observability

* Start one **OTel span per repository method** (e.g., `repo.category.create`, `repo.category.get_by_id`, …).
* Attach canonical attributes:

    * `operation`, `user_id`, `category_id`, `status`.
* On error: `span.RecordError(err)` + set status; log **metadata** only (no sensitive payloads).

---

## Mapping rules

* **IDs & ownership**: user ownership is enforced at the query level (filter by `user_id`).
* **Timestamps**: `updated_at` is set on writes; `deleted_at` is set by **soft delete** (no hard deletes).
* **Nullables**: optional fields (`description`, `color_hex`, `icon`) map 1:1 between DB and domain.

---

## Design notes

* This adapter **must not** leak ORM types to the core; always map to `domain.Category`.
* Prefer **explicit columns** and **scoped queries** (id + user\_id / name + user\_id).
* Keep **update** APIs partial: accept a `map[string]any` of changed fields only.
* Soft delete only marks `deleted_at`; readers should not return deleted items.

---

## Testing hints

* Unit tests for the **core** should mock the output port—no DB.
* Adapter tests can use:

    * a disposable Postgres (Docker) or
    * GORM against an ephemeral DB.
* Assert:

    * correct filtering by `user_id`
    * mapping DB ⇄ domain
    * soft delete behavior (row persists, `deleted_at` set)
    * error propagation and span logging (can be observed via test logger).

---

## Performance & operational tips

* Queries typically filter by `(user_id, name)` or `(user_id, id)` — ensure appropriate DB indexes if needed.
* Keep migrations aligned with the model tags and table names (`aion_api.tag_categories`).
* Use transactions in the repository layer if a future operation spans multiple statements.

---

## Gotchas

* Don’t return soft‐deleted rows from read queries.
* Keep **domain** free from GORM tags/types; all ORM specifics belong to `model/`.
* Always pass the **request context** into GORM (`db.WithContext(ctx)`) to honor timeouts/cancellation.
