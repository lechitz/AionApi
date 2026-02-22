# Documentation Ownership (Canonical Map)

Last verified: 2026-02-22

Purpose: define canonical documentation sources and reduce doc drift.

## Canonical Sources
- REST contract: `contracts/openapi/swagger.yaml`
- GraphQL contract set: `contracts/graphql/`
- GraphQL schema artifact: `docs/graphql/schema.graphql`
- Architecture overview: `README.md` and `docs/architecture.md`
- Runtime/infrastructure ops: `infrastructure/README.md`
- AI governance rules: `AGENTS.md` and `.github/PULL_REQUEST_TEMPLATE.md`

## Rules
- If docs conflict, contracts and canonical files above win.
- Contract changes must update canonical docs in the same PR.
- Narrative docs must link to canonical contracts when applicable.
