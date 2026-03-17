# AionAPI - Excellence Guide for Agents

If this repo is being used from the workspace root `../`, follow the workspace root `AGENTS.md` for global shortcuts, `!issue` behavior, and Project 2 operations. This file remains the local authority for AionAPI architecture and implementation rules.

This document defines the professional standard expected when working on AionAPI. Follow these rules for any change, review, or proposal.

## 1. Project purpose
AionAPI is a professional Go backend built with Hexagonal (Ports & Adapters) and Clean Architecture, organized by bounded contexts.

## 2. Non-negotiable rules
- **Dependency rule:** `internal/<ctx>/core/{domain,ports,usecase}` must not import adapters, platform, or infrastructure. `domain` is pure Go (stdlib only).
- **Ports as contracts:** Input/output ports live in `internal/<ctx>/core/ports/input` and `internal/<ctx>/core/ports/output`. Usecases implement input ports and depend only on output ports.
- **Context isolation:** No direct imports between bounded contexts. Share only via `internal/shared` or explicit events/ports.
- **Thin adapters:** Primary adapters translate transport <-> core. Secondary adapters implement output ports and translate infra errors to semantic errors.
- **Central GraphQL:** Resolvers live in `internal/adapter/primary/graphql` and delegate to context controllers in `internal/<ctx>/adapter/primary/graphql/controller`.
- **Error model:** Use semantic errors from `internal/shared/sharederrors`. Adapters map to HTTP/GraphQL responses without leaking infra details.
- **Observability:** Every handler/controller/usecase/repository opens an OTel span and propagates context.

## 3. Code standards
- **Idiomatic Go:** small functions, explicit errors, clear names, no hidden side effects.
- **Context required:** `context.Context` at all boundaries (adapters <-> usecases <-> repositories).
- **Structured logs:** use `logger.ContextLogger`; log metadata, never sensitive data.
- **No magic strings:** shared constants in `internal/shared/constants`; local constants in `*_constants.go` within the package.
- **Mapping discipline:** keep DTO/GraphQL model mapping in adapters/controllers, never in usecases.

## 4. Flow responsibilities
- **Primary adapters (HTTP/GraphQL):** validation, tracing, DTO <-> domain mapping, call usecase, map response.
- **Usecases:** business orchestration, no infrastructure dependencies.
- **Secondary adapters:** persistence/cache/external IO, always context-aware, translate infra errors -> semantic.

## 5. Testing and quality
- **Usecases:** unit tests with `gomock` and table-driven cases.
- **Adapters:** test mappings and boundary behavior (GraphQL controllers, HTTP handlers).
- **Coverage and mocks:** when adding ports, generate mocks and cover new flows.
- **Official commands:** `make mocks`, `make graphql`, `make test`, `make test-cover-detail`, `make verify`.

## 6. Codegen and schema changes
- **GraphQL:** update `.graphqls` under `internal/adapter/primary/graphql/schema/modules/`, then run `make graphql`.
- **Mocks:** after adding/changing output ports, run `make mocks` (mocks live in `tests/mocks`).
- **Migrations:** add SQL migrations under `infrastructure/db/migrations` and keep them reversible.

## 7. AI control center (local)
- The `/agents` folder is the AI control center (Codex, Copilot, etc.). Everything must follow what is defined there.
- Always select the correct persona from `/agents/personas` before proposing changes.
- If `/agents` does not exist in the current environment, follow this `AGENTS.md` as the default.
- Never cite private paths or content in public docs or PRs.

## 8. Review checklist
- Architecture boundaries respected (no cross-context imports).
- Errors are semantic and mapped in adapters.
- Spans/logs added at boundaries with safe metadata only.
- Mapping kept in adapters/controllers; usecases remain pure.
- Tests updated for new behavior and ports.

## 9. Decision protocol
- Select the persona that best matches the task. If multiple apply, follow the strictest architectural rule first.
- When rules conflict, default to architecture boundaries and error model, then observability, then style.
- For ambiguous placement ("where does this go?"), mirror the `category` context structure and confirm paths.

## 10. Change discipline
- Every proposal must state impact on architecture, tests, and observability (even if "no impact").
- Avoid cross-context changes unless explicitly required; prefer local refactors with clear scope.
- Large refactors require a short plan before edits.

## 11. Definition of done (DoD)
- Correct layer placement and dependency direction.
- Spans and structured logs at boundaries.
- Semantic errors used and mapped in adapters.
- Tests or rationale for why tests are not required.

## 12. Risk flags (stop and reassess)
- Cross-context imports or shared state leakage.
- Business logic in adapters or transport layers.
- Infra types in domain or usecases.
- Unhandled errors or generic `error` returns where a semantic error exists.

## 13. Expected deliverables
- Clear file paths and objective rationale for changes.
- Proposals aligned with existing high-quality contexts (e.g., `category`).
- Avoid large refactors without need; prefer small, consistent changes.

## 14. Local shortcut note
- Workspace shortcuts such as `!issue`, `!bug`, `!review`, and `!ship` are defined globally in `../AGENTS.md`.
- For repo-local implementation, apply those global shortcuts with AionAPI's architectural constraints and the personas in `/agents`.
