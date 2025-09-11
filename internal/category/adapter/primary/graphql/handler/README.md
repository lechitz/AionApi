# Category GraphQL Controller (Primary Adapter)

**Path:** `internal/category/adapter/primary/graphql/handler`

## Responsibilities

* Provide **GraphQL-facing controllers** for the Category context.
* Map **GraphQL models ⇄ domain entities**, call **CategoryService** (input port), and return GraphQL-friendly results.
* Orchestrate **tracing** (OTel) and **structured logging**; no business rules here.

## How it works

* **Dependency Injection:** `NewHandler(categoryService, logger)` receives the input port and logger.
* **Controller methods (called by central resolvers):**

    * `Create(ctx, input, userID) (*model.Category, error)`
    * `Update(ctx, input, userID) (*model.Category, error)`
    * `SoftDelete(ctx, categoryID, userID) error`
    * `ListAll(ctx, userID) ([]*model.Category, error)`
    * `GetByID(ctx, categoryID, userID) (*model.Category, error)`
    * `GetByName(ctx, name, userID) (*model.Category, error)`
* **Mapping helpers:** `toDomainCreate`, `toDomainUpdate`, `toModelOut` isolate model/domain translation.
* **User identity:** `userID` must be provided by the **upper layer** (e.g., GraphQL resolver/middleware). This controller doesn’t extract auth.

## Why this package matters

* Keeps **GraphQL transport concerns** (shape/mapping/tracing/logging) **out of the core**.
* Maintains a **thin, testable** seam between GraphQL resolvers and domain use cases.
* Enables evolution of schema/transport without touching business logic.

## Observability

* Tracer name: `aionapi.graphql.handler`.
* Spans per operation (e.g., `handler.create`, `handler.update`, `handler.list_all`).
* Adds canonical attributes (e.g., `user_id`, `category_id`, `category_name`) and records errors.

## Error handling

* Validates presence of `userID` and basic input (e.g., parse IDs).
* Returns **semantic errors** upward; does not shape HTTP/GraphQL responses here (that’s resolver/transport’s job).

## Do / Don’t

**Do**

* Keep handlers **thin**: validate input, map, call `CategoryService`, map output, add trace/logs.
* Propagate `context.Context` for deadlines/cancellation/tracing.
* Log **metadata only** (never secrets); use canonical keys.

**Don’t**

* Implement business rules (belongs to `core/usecase`).
* Access DB/cache directly (use the input port).
* Perform authentication/authorization here (ensure `userID` is already resolved by middleware/resolver).
