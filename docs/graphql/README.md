# GraphQL Documentation Artifacts

**Path:** `docs/graphql`

## Overview

This folder contains generated GraphQL documentation artifacts used by consumers and tooling.
It complements schema sources in `internal/adapter/primary/graphql/schema` and shared operations in `contracts/graphql/queries`.

## Files

| File | Purpose |
| --- | --- |
| `schema.graphql` | Flattened SDL representation of current GraphQL schema |
| `introspection.json` | GraphQL introspection artifact for client/tooling usage |

## Regeneration Workflow

```bash
make graphql.schema
make graphql.introspect
```

## Related Sources

- Schema modules: `internal/adapter/primary/graphql/schema/modules/`
- Shared queries: `contracts/graphql/queries/`

## Design Notes

- Treat artifacts here as generated outputs.
- Keep docs artifacts in sync with schema contract changes.
- Use these files for client/codegen tooling integration.

## Package Improvements

- Add CI guard for stale generated schema/introspection artifacts.
- Add explicit generation prerequisites (running server for introspection).
- Add compatibility note for tooling versions that consume introspection.
- Add quick links to GraphQL playground and central adapter README.

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
