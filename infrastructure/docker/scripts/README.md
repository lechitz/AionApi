# infrastructure/docker/scripts

Shell scripts used by Docker workflows, Make targets, and CI to build/run/clean the environment.

## Package Composition

- `entrypoint.sh`
  - Container entrypoint for the API image.
- `README.md`
  - Script conventions and usage notes.

## Flow (Where it comes from -> Where it goes)

Make/CI -> scripts -> docker build/run lifecycle

## Why It Was Designed This Way

- Keep Docker concerns isolated from application code.
- Centralize environment initialization in one place.

## Recommended Practices Visible Here

- Prefer `.env.*` over hardcoded values.
- Keep scripts portable (sh/bash) and documented.
- Make destructive actions explicit and opt-in.

## What Should NOT Live Here

- Business logic or app behavior.
- Secrets embedded in scripts.
