# docs/swagger

Generated OpenAPI/Swagger artifacts for the AionAPI REST contract. These files feed Swagger UI and are published for client consumption.

## Package Composition

- `swagger.yaml` / `swagger.json`
  - Generated REST contract.
- `docs.go`
  - Swaggo binding used by the API binary.

## Flow (Where it comes from -> Where it goes)

Handler annotations -> swagger.{yaml,json} -> Swagger UI / client tooling

## Why It Was Designed This Way

- Keep the REST contract generated from source code.
- Avoid manual drift between code and docs.

## Recommended Practices Visible Here

- Regenerate via Make targets; do not edit generated files.
- Keep error envelopes aligned with `internal/platform/server/http/utils`.
- Version breaking changes explicitly.

## Differentials (Rare but Valuable)

- Contract artifacts are versioned and published.

## How to Generate

```bash
make docs.gen
make docs.check-dirty
make docs.clean
```

## What Should NOT Live Here

- Manual edits of generated files.
- Source-of-truth API definitions.
