# Shared GraphQL Queries

**Path:** `contracts/graphql/queries`

## Overview

This folder stores reusable GraphQL operation documents shared across Aion clients.
It acts as a versioned contract artifact aligned with the backend GraphQL schema.

## Package Scope

| Area | Responsibility |
| --- | --- |
| Shared operations | Keep canonical query documents in one place |
| Client reuse | Allow multiple clients to load the same `.graphql` files |
| Contract traceability | Version query changes together with backend schema evolution |

## Structure

| Folder | Current operations |
| --- | --- |
| `categories/` | `list.graphql` |
| `records/` | `list.graphql` |
| `tags/` | `list.graphql` |

## Current Query Inventory

| Domain | File | Operation |
| --- | --- | --- |
| Categories | `categories/list.graphql` | `ListCategories` |
| Records | `records/list.graphql` | `ListRecords(limit: Int)` |
| Tags | `tags/list.graphql` | `ListTags` |

**Total operations:** `3`

## Usage Examples

### Python

```python
from pathlib import Path

query = (Path("contracts/graphql/queries/categories/list.graphql")).read_text()
```

### TypeScript

```ts
import { readFileSync } from "fs";

const query = readFileSync("contracts/graphql/queries/categories/list.graphql", "utf-8");
```

## Maintenance Workflow

1. Add or update `.graphql` files under the appropriate domain folder.
2. Ensure operation names and selected fields match current schema.
3. Validate client compatibility (dashboard, chat, CLI, etc.).
4. Commit query files together with related schema/backend changes.

## Design Notes

- Keep operation documents small and purpose-driven.
- Prefer stable operation names to improve observability and cache behavior.
- This folder stores operation documents only, not generated clients.

## Package Improvements

- Add missing mutation/query documents for create/update/delete flows to match schema coverage.
- Add CI validation that checks `.graphql` files against the current GraphQL schema.
- Standardize folder naming (`categories`, `tags`, `records`) with a short naming convention note.
- Add per-query comments for expected auth context (`@auth` protected operations).

---

<!-- doc-nav:start -->
## Navigation
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
