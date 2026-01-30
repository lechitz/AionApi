# internal/adapter/primary/graphql/directives

Shared GraphQL directives that enforce cross-cutting policies (auth, roles, validation, tracing).

## Package Composition

- Directive implementations and helpers.

## Flow (Where it comes from -> Where it goes)

Schema directive -> directive handler -> resolver/controller -> usecase

## Why It Was Designed This Way

- Apply policies consistently at the schema level.
- Keep resolvers focused on domain orchestration.
- Provide uniform error responses for policy failures.

## Recommended Practices Visible Here

- Propagate context values needed for observability.
- Keep directive errors consistent for clients.
- Document behavior when adding new directives.

## Differentials

- Policy enforcement at schema level instead of resolver duplication.

## What Should NOT Live Here

- Business rules or persistence logic.
- Context-specific validations that belong to usecases.
