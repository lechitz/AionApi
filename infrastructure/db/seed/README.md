# Database Seed Scripts

**Path:** `infrastructure/db/seed`

## Overview

SQL seed scripts used for local development, QA datasets, and deterministic test baselines.
These scripts populate data directly in Postgres for fast environment setup.

## Scope

| Area | Responsibility |
| --- | --- |
| Core seed data | Roles, users, relationships |
| Generated datasets | Categories, tags, records, test timelines |
| Seed logging/artifacts | Optional run artifacts for local tooling |

## Key Files

| File pattern | Purpose |
| --- | --- |
| `*_generate.sql` | Parameterized data generators |
| `roles.sql`, `user_roles.sql`, `admin_user.sql` | Security/role bootstrap |
| `test_*.sql` | Scenario-focused test datasets |

## Usage

```bash
make seed-all N=10
make populate N=100
make db-full
```

## Realistic Demo Flow

Use this flow when you want a production-like local environment:

```bash
make dev
make db-full
```

What `db-full` prepares:
- reset + full migrations
- roles + admin account
- realistic user `testuser` / `Test@123`
- merged taxonomy (legacy + new operational categories/tags)
- dashboard metric definitions + goal templates
- dashboard white-label layout tables ready (`dashboard_views`, `dashboard_widgets`)
- ~3 months of records with high daily volume (about 50-60/day)

## Design Notes

- Prefer idempotent SQL where possible.
- Keep data representative but lightweight for local workflows.
- Avoid embedding secrets in seed files.

## Package Improvements

- Add automated seed validation against current schema version.
- Add documented dependency order between seed scripts.
- Add optional “small/medium/large” seed profiles.
- Move local artifacts (e.g., transient logs) to ignored output folder.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
