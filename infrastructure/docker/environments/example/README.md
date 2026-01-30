# infrastructure/docker/environments/example

Template environment for Docker overlays and `.env.example`. Use this as the base for creating new environments.

## Package Composition

- `.env.example`
  - Full list of required variables (HTTP, GraphQL, DB, Cache, OTEL, JWT).
- `docker-compose-*`
  - Starting point for new profiles (copy and adjust ports/endpoints).

## Flow (Where it comes from -> Where it goes)

Template -> new environment -> docker compose + app config

## Why It Was Designed This Way

- Provide a single source of truth for required env vars.
- Reduce drift across environment definitions.

## Recommended Practices Visible Here

- Keep the template in sync when new variables are added.
- Document expected formats in the file (URLs, durations, secrets).
- Validate with `make dev` or the target for the new profile.

## What Should NOT Live Here

- Real secrets or production values.
- Environment-specific overrides (create a new env folder).
