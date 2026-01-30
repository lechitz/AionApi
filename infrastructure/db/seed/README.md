# infrastructure/db/seed

Seed data scripts for local development and manual QA. These files populate a consistent dataset without calling the API directly.

## Package Composition

- `*_generate.sql`
  - Parametrized generators for users, categories, tags, and records.
- `roles.sql`, `admin_user.sql`, `user_roles.sql`
  - System roles and admin/user role assignments.

## Flow (Where it comes from -> Where it goes)

Developer -> make seed-* -> SQL scripts -> database tables

## Why It Was Designed This Way

- Provide fast, deterministic datasets for dev/test.
- Keep seeding reproducible and idempotent.
- Avoid API dependency when validating DB structure.

## Recommended Practices Visible Here

- Keep scripts idempotent (ON CONFLICT DO NOTHING).
- Preserve foreign key order in execution.
- Use small N for quick local runs and larger N for load tests.

## Differentials

- Parametrized generators for scalable dataset sizes.
- Consistent naming to avoid collisions across runs.

## Common Targets

```bash
make seed-all N=10
make populate N=100
make seed-user1-all
```

## Environment Variables

| Variable | Default | Description |
| --- | --- | --- |
| `N` | 10 | Number of users to generate |
| `DEV_PASSWORD` | testpassword123 | Password for all seeded users |
| `SEED_DAYS` | 7 | Number of days of records to generate |

## What Should NOT Live Here

- Production credentials or secrets.
- API-driven seed logic (use `cmd/api-seed-caller`).
- Irreversible or destructive SQL without rollback strategy.
