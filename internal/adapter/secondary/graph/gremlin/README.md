# internal/adapter/secondary/graph/gremlin

Gremlin driver integration for graph ports.

## Package Composition

- Gremlin client configuration and traversal helpers.

## Flow (Where it comes from -> Where it goes)

Usecase -> graph port -> Gremlin adapter -> graph database

## Why It Was Designed This Way

- Centralize Gremlin connection and traversal patterns.
- Keep vendor-specific details out of core logic.

## Recommended Practices Visible Here

- Document any vendor-specific traversal extensions.
- Use timeouts to avoid long-running traversals.
- Capture latency and error metrics for diagnostics.

## Differentials

- Reusable Gremlin traversal helpers.

## What Should NOT Live Here

- Business rules or domain logic.
