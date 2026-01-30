# internal/tag

Tag domain for user-scoped labels exposed via GraphQL with per-user uniqueness.

## Purpose and Main Capabilities

- Create and query tags by ID/name/category.
- Enforce required fields and per-user uniqueness.
- Keep tag validation centralized in core.
- Provide semantic errors for GraphQL mapping.

## Package Composition

- `core/`: tag entities, ports, and usecases.
- `core/ports/input`: tag service interface for adapters.
- `core/ports/output`: tag repository contract.
- `core/usecase`: Create, GetByID, GetByName, GetAll, GetByCategory.
- `adapter/primary/graphql`: GraphQL controllers/resolvers and DTO mapping.
- `adapter/secondary/db`: persistence adapters and mappers.

## Flow (Where it comes from -> Where it goes)

GraphQL request -> context controller -> input port -> usecase ->
output port -> db adapter -> database -> response

## How It Works (Concise)

- Controllers map GraphQL inputs and read `user_id` from context.
- Usecases validate name/description and check uniqueness before create.
- Repositories persist and query tags by user and category.

## Why It Was Designed This Way

- Keep tag rules consistent and transport-agnostic.
- Avoid duplication of validation across resolvers.
- Make schema and persistence evolution straightforward.

## Recommended Practices Visible Here

- Validate in core and return semantic errors (`TagAlreadyExists`, `TagNameIsRequired`).
- Use spans with `tag_name`, `category_id`, `user_id` only.
- Keep mapping in adapters; core stays pure.

## Differentials

- Strict per-user uniqueness baked into the usecases.

## What Should NOT Live Here

- Business logic inside adapters.
- Transport DTOs inside core.
- Cross-context imports.
