# infrastructure

Operational assets for AionAPI: build, runtime environments, data lifecycle, and observability stack.

## Package Composition

- `docker/`
  - Dockerfiles, compose overlays, and environment templates.
- `db/`
  - Migrations and seed data for Postgres.
- `observability/`
  - Logs, metrics, and tracing infrastructure.

## Flow (Where it comes from -> Where it goes)

Developer config -> infrastructure assets -> local/prod runtime stack

## Why It Was Designed This Way

- Keep operational contracts versioned alongside code.
- Make local setups reproducible.
- Separate runtime concerns (docker, db, observability) cleanly.

## Recommended Practices Visible Here

- Keep environment differences isolated under `docker/environments/`.
- Align infra endpoints with app config and docs.
- Treat migrations as immutable, seeds as deterministic.

## Differentials

- Full local stack documented and versioned.
- Clear separation between runtime wiring and app code.

## What Should NOT Live Here

- Application business logic.
- Secrets or environment-specific credentials.
- One-off scripts without repeatable value.
