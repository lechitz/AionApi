# internal/category

Category taxonomy domain (habit/daily) exposed via GraphQL, with per-user uniqueness and soft delete.

## Purpose and Main Capabilities

- Provide CRUD for categories scoped to a user.
- Enforce uniqueness of category name per user.
- Support partial updates and soft delete semantics.
- Keep category rules isolated from transport and persistence.

## Package Composition

- `core/`: domain rules, ports, and usecases.
- `core/ports/input`: Category service interface for adapters.
- `core/ports/output`: Category repository contract.
- `core/usecase`: Create, GetByID, GetByName, ListAll, Update, SoftDelete.
- `adapter/primary/graphql`: GraphQL controllers/resolvers and DTO mapping.
- `adapter/secondary/db`: repository, mapper, and persistence models.

## Flow (Where it comes from -> Where it goes)

GraphQL request -> context controller -> input port -> usecase ->
output port -> db adapter -> database -> response

## How It Works (Concise)

- Controllers read `user_id` from context and map GraphQL inputs to commands.
- Usecases validate invariants (name required, per-user uniqueness) and orchestrate repository calls.
- Updates are partial to avoid unintended overwrites.
- Soft delete is idempotent to preserve referential integrity with records.

## Why It Was Designed This Way

- Keep category rules consistent and reusable across transports.
- Avoid leaking ORM/persistence types into core.
- Make category behavior safe for auditing and future extension.

## Recommended Practices Visible Here

- Validate uniqueness in core before persist.
- Keep mapping in adapters; core stays pure.
- Emit spans with `user_id`/`category_id` and log only metadata.
- Coordinate schema and persistence changes (GraphQL + migrations).

## Differentials

- Strict per-user uniqueness and soft delete by design.

## What Should NOT Live Here

- Cross-context imports or shared state.
- Transport DTOs inside core.
- Business logic inside adapters.
