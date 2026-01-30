# infrastructure/docker/environments

Docker compose overlays and `.env` files for each runtime profile. This folder defines how the stack runs in dev and prod.

## Package Composition

- `dev/`
  - Local development profile and `.env.dev`.
- `prod/`
  - Production profile and `.env.prod`.
- `example/`
  - Template for new environments.

## Flow (Where it comes from -> Where it goes)

Environment profile -> docker compose -> running services

## Why It Was Designed This Way

- Separate dev and prod runtime settings.
- Keep secrets isolated per environment.
- Provide a template for consistent new profiles.

## Recommended Practices Visible Here

- Always derive new environments from `example/`.
- Keep env files per profile; never reuse dev env in prod.
- Update Makefile targets when adding a new profile.

## What Should NOT Live Here

- Real secrets committed to git.
- One-off local overrides (use a local copy).
