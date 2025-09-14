# Category Usecases (Core)

**Path:** `internal/category/core/usecase`

## Responsibilities

* Implement domain rules for the **Category** context behind input ports.
* Validate inputs, enforce invariants/policies, and orchestrate calls to output ports (repository).
* Normalize persistence/infra errors into **semantic domain errors**.

## How it works

* **Construction:** `NewService(repository output.CategoryRepository, logger logger.ContextLogger) *Service`
* **DIP:** Depends only on **interfaces** from `core/ports/output`; concrete adapters are wired in `platform/bootstrap`.
* **Public API (examples):**

    * `Create(ctx, category)` → validates required fields, checks **uniqueness by name** (`GetByName`), persists via repo, returns created entity.
    * `GetByID(ctx, categoryID, userID)` → fetches by composite key; validates `categoryID`.
    * `GetByName(ctx, name, userID)` → validates `name`, fetches via repo.
    * `ListAll(ctx, userID)` → reads all categories for a user.
    * `Update(ctx, category)` → builds **partial update map** from non-empty fields and delegates to repo.
    * `SoftDelete(ctx, categoryID, userID)` → delegates soft deletion; idempotent by design.
* **Validation & errors:** Guard clauses for required fields and limits (name/description/color/icon). Returns meaningful constants (e.g., `CategoryNameIsRequired`, `FailedToCreateCategory`).
* **Observability:** Opens OTel spans (`TracerName`, `SpanCreateCategory`, etc.), sets canonical attributes (user/category IDs), emits events (`EventRepositoryCreate`, …), records errors; uses structured logging.

## Typical flow

1. Handler (HTTP/GraphQL) calls an **input port** on `Service` with context + values.
2. Usecase validates inputs/invariants and normalizes data.
3. Calls repository via **output port**; handles not-found/DB errors consistently.
4. Returns **domain entities** or **semantic errors** to the adapter.

## My reminders

* Keep usecases **transport-agnostic** (no HTTP/GraphQL/ORM types).
* Favor small, **pure** functions; easy to test with mocks.
* **Normalize** repository “not found” and DB errors into domain errors.
* Build update maps explicitly; never allow unintended/immutable fields.
* Always pass `context.Context`; set meaningful span attributes and events.
* Log **metadata only** (IDs, counts), never sensitive payloads.
* Cover edge cases with unit tests (validation, uniqueness, partial updates, soft delete).
