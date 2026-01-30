# docs/collections

This folder contains versioned API client collections (Postman/Insomnia) for local development and manual QA. OpenAPI remains the source of truth.

## Package Composition

- `postman/`
  - Current collection file(s), for example `AionApi.postman_collection.json`.

## Flow (Where it comes from -> Where it goes)

OpenAPI spec -> Client collection -> Manual QA / local testing

## Why It Was Designed This Way

- Keep consumer tooling in sync with the API contract.
- Support quick manual testing without custom scripts.

## Recommended Practices Visible Here

- Keep `{{baseURL}}` and environment variables.
- Version collections on breaking changes.
- Avoid committing secrets or tokens.

## Differentials (Rare but Valuable)

- Collections are versioned artifacts, not the contract source.

## What Should NOT Live Here

- API contract changes (use `docs/swagger`).
- Secrets or real tokens.
