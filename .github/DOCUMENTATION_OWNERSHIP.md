# Documentation Ownership (Canonical Map)

Last verified: 2026-03-17

Purpose: define canonical documentation sources and reduce doc drift.

## Canonical Sources
- REST contract: `contracts/openapi/swagger.yaml`
- GraphQL schema modules + resolver behavior: `internal/adapter/primary/graphql/schema/modules/` and `internal/record/adapter/primary/graphql/controller/`
- GraphQL shared contract set: `contracts/graphql/`
- GraphQL manifest: `contracts/graphql/manifest.json`
- GraphQL flattened schema artifact: `docs/graphql/schema.graphql`
- GraphQL consumer-facing narrative docs: `docs/graphql/README.md` and `docs/graphql-playground.md`
- Architecture overview: `README.md` and `docs/architecture.md`
- Runtime/infrastructure ops: `infrastructure/README.md`
- AI governance rules: `AGENTS.md` and `.github/PULL_REQUEST_TEMPLATE.md`
- Cross-repo contract governance: `aion-docs/planning/v1/adr/adr-005-insight-contract-policy.md`

## Rules
- If docs conflict, contracts and canonical files above win.
- For GraphQL, authority order is:
  1. live schema modules and backend behavior
  2. generated backend artifacts
  3. shared contract artifacts and manifest
  4. narrative docs and playground pages
- GraphQL contract changes must update the relevant canonical artifacts in the same PR:
  - schema modules and resolver/controller behavior
  - `make graphql.schema`
  - `make graphql.queries graphql.manifest`
- Narrative docs must link back to canonical contracts instead of redefining semantics.
- Generated artifacts under `docs/graphql/` are snapshots and should not be hand-edited to describe behavior that is owned elsewhere.
