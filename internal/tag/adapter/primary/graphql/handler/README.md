# Tag GraphQL Controller (Primary Adapter)

**Folder:** `internal/tag/adapter/primary/graphql/handler`

## Responsibilities

* Orchestrate Tag operations exposed via GraphQL resolvers.
* Map **GraphQL inputs ⇄ domain entities** and return GraphQL models.
* Add tracing/log fields (consistent with the Category controller).
* Surface domain/usecase errors without embedding transport/ORM details.

## How it works

* `NewHandler(svc input.TagService, logger logger.ContextLogger)` injects the **TagService** (input port) and logger.
* Public methods are called **by the central GraphQL resolvers** (`internal/adapter/primary/graph/resolver/*`):

    * `Create(ctx, in model.CreateTagInput) (*model.Tag, error)`
    * `GetByID(ctx, id string) (*model.Tag, error)`
    * `ListAll(ctx) ([]*model.Tag, error)`
* Each method should:

    1. Start an OTel span (use the same naming style as Category).
    2. Map GraphQL `model.*` → domain types.
    3. Call the **TagService** (core/usecase via input ports).
    4. Map domain → GraphQL `model.*` and return.
    5. Log with canonical keys (user\_id, tag\_id, operation, etc.).

## Reminders

* **No business logic here**—delegate to `core/usecase`.
* **Do not access DB**—only call ports on the service.
* Keep mapping helpers (`toDomainCreate`, `toModelOut`, etc.) in this package.
* Use `context.Context` for tracing/metadata; extract `user_id` from context if tags are user-scoped.
* Return **semantic errors** from the service as-is; avoid ad-hoc string errors.

## Current status / TODO

* Methods are **stubbed**; implement them mirroring the Category controller pattern:

    * Add spans, attributes, and structured logs.
    * Implement mapping helpers.
    * Decide whether tags are **global or per user**; if per user, require/propagate `userID`.
* Wire resolvers to call these methods (if not already).
* Add unit tests for mapping and error propagation.
