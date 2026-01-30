# internal/adapter/secondary/graph

Graph database adapters compatible with Gremlin-style APIs.

## Package Composition

- `gremlin/`
  - Gremlin client implementation and traversal helpers.

## Flow (Where it comes from -> Where it goes)

Usecase -> graph port -> graph adapter -> graph database

## Why It Was Designed This Way

- Encapsulate traversal construction and driver specifics.
- Keep graph queries swappable across vendors.

## Recommended Practices Visible Here

- Keep traversals composable and documented.
- Apply timeouts/retries for graph calls.
- Add observability around graph operations.

## Differentials

- Driver-specific logic isolated from core.

## What Should NOT Live Here

- Domain rules or transport logic.
