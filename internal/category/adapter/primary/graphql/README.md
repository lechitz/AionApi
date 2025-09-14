# Category GraphQL (Primary Adapter)

**Folder:** `internal/category/adapter/primary/graphql`

## Responsibility

* Expose the **Category** context over **GraphQL** (handlers/controllers, resolvers, schema).
* Keep controllers **thin**: map GraphQL models ⇄ domain, call input ports, and standardize errors.
* Orchestrate **observability** (OpenTelemetry spans + structured logs). No business rules here.

## How it works

* **Structure**

    * `handler/` — GraphQL controllers that:

        * start spans (`TracerName = "aionapi.graphql.handler"`)
        * read context (e.g., `user_id`) set by upstream auth
        * map inputs to domain and call `CategoryService` (input port)
    * `resolver/` — gqlgen resolvers that delegate to the handler.
    * `schema/` — GraphQL schema pieces for the Category domain (`*.graphqls`).

* **DI & flow**

    * `NewHandler(categoryService, logger)` injects dependencies.
    * Resolvers receive a `*handler.Handler` and call its methods (e.g., `Create`, `Update`, …).
    * Handlers call the **input port** (`internal/category/core/ports/input`) — never talk to repositories/ORM.

## Operations (typical)

> Exact names come from `schema/category.graphqls`. The handler currently supports:

**Queries**

* Get by ID
* Get by name
* List all (for the authenticated user)

**Mutations**

* Create category
* Update category (partial updates)
* Soft delete category

> Handlers read `user_id` from the request context (populated by the upstream auth layer).

## Controller conventions

* Start one span per operation:

    * Spans (examples): `handler.create`, `handler.update`, `handler.soft_delete`, `handler.list_all`, `handler.get_by_id`, `handler.get_by_name`.
    * Set attributes like `operation`, `user_id`, and status labels (e.g., `category_created`, `category_updated`).
* Use **only** the input port (`CategoryService`) — keep transport concerns (GraphQL I/O, mapping, tracing) in this layer.
* Map GraphQL inputs to domain structs; **do not** expose domain types directly in the API.
* Return GraphQL-friendly results and standardized errors.

## Schema & codegen

* Category schema lives in `schema/category.graphqls`.
* Run `make graphql` after schema changes:

    * The build copies all `*.graphqls` into the gqlgen module and regenerates models/resolvers.
    * Always `go mod tidy` after generation (handled by the make target).

## Observability

* **Tracer:** `aionapi.graphql.handler`
* **Spans:** see constants in `handler/0_category_handler_constants.go`.
* **Best practices:**

    * Record errors on the span (`span.RecordError(err)`) and set status codes.
    * Log **metadata** (IDs, operation, status), never sensitive payloads.

## Testing hints

* Unit-test handlers by stubbing `CategoryService`:

    * Success paths (create/update/list/get) → assert mapping and returned models.
    * Error paths (service errors/not found/validation) → assert span error + standardized error mapping.
* For resolvers, use a fake handler and assert that GraphQL inputs are forwarded/mapped correctly.
* Relax logger expectations in gomock (e.g., `.AnyTimes()`) to keep tests focused on behavior.

## Design notes

* Transport-agnostic business logic lives in **use cases**; this adapter is **I/O + orchestration** only.
* Keep resolvers minimal — they should forward to the handler and return mapped results.
* Keep auth outside the handler; expect `user_id` in context (set by the platform/auth layer).
