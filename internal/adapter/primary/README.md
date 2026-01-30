# internal/adapter/primary

Primary adapters expose domain usecases to external clients (GraphQL/HTTP). They translate transport payloads into domain commands and map responses back.

## Package Composition

- `graphql/`
  - Shared GraphQL schema, directives, and models.

## Flow (Where it comes from -> Where it goes)

Client request -> primary adapter -> controller/usecase -> response DTO

## Why It Was Designed This Way

- Keep transport concerns out of core logic.
- Provide consistent auth/validation and error mapping.
- Reuse shared GraphQL infrastructure across contexts.

## Recommended Practices Visible Here

- Handlers/resolvers stay thin; business rules live in core.
- Reuse directives/middleware for cross-cutting policies.
- Avoid leaking transport types into domain.

## Differentials

- Centralized transport conventions to reduce drift.

## What Should NOT Live Here

- Business logic or persistence.
- Domain entities with transport tags.
