# Primary GraphQL Adapter

**Path:** `internal/adapter/primary/graphql`

## Overview

This package is the central GraphQL transport entrypoint for the application.
It hosts schema composition, gqlgen-generated artifacts, directive wiring, and thin resolvers that delegate to context controllers.

## Package Scope

| Area | Responsibility |
| --- | --- |
| Server setup | Build GraphQL HTTP handler and configure transports/middlewares |
| Schema composition | Maintain root schema and module-based schema extensions |
| Directive wiring | Register cross-cutting directives (currently `@auth`) |
| Resolver bridge | Convert GraphQL inputs/context into controller calls |
| Codegen integration | Define gqlgen generation targets/configuration |

## Structure

| Path | Role |
| --- | --- |
| `schema/root.graphqls` | Root Query/Mutation and shared directives/scalars |
| `schema/modules/*.graphqls` | Domain extensions (`category`, `tags`, `record`, `chat`, `user`) |
| `directives/auth.go` | Authorization directive implementation (`@auth`) |
| `resolver.go` | Dependency wiring from services to context GraphQL controllers |
| `*.resolvers.go` | Thin field resolvers generated/preserved by gqlgen |
| `server.go` | HTTP router, recovery middleware, auth middleware, transports |
| `generated.go` | gqlgen generated execution engine (do not edit manually) |
| `model/models_gen.go` | gqlgen generated GraphQL transport models |
| `gqlgen.yml` | gqlgen configuration (schema, resolver/model outputs) |

## GraphQL Runtime Flow

1. HTTP request enters GraphQL handler (`server.go`).
2. Recovery middleware wraps request execution.
3. Auth middleware populates request context (when auth service is configured).
4. gqlgen executes operation and runs directives (e.g., `@auth`).
5. Resolver reads context values (e.g., `ctxkeys.UserID`) and delegates to context controller.
6. Controller maps GraphQL DTOs to usecase calls and returns response data/errors.

## Directive Model

| Directive | Behavior |
| --- | --- |
| `@auth(roles: String)` | Requires `ctxkeys.UserID`; enforces role when provided; bypasses role checks for trusted `ctxkeys.ServiceAccount=true` calls |

## Configured Transports

| Transport | Status |
| --- | --- |
| `GET` | Enabled |
| `POST` | Enabled |
| `OPTIONS` | Enabled |
| `MultipartForm` | Enabled |
| `Websocket` | Enabled (`KeepAlivePingInterval: 10s`) |

## Design Notes

- Resolvers must remain thin and should not host domain business logic.
- Context-specific mapping/orchestration belongs to `internal/<ctx>/adapter/primary/graphql/controller`.
- Generated files (`generated.go`, `model/models_gen.go`) should be managed by codegen only.
- Module-based schema (`schema/modules`) keeps domain contracts isolated while exposing a single endpoint.

## Package Improvements

- Add resolver tests for parsing/validation edges (e.g., invalid ID conversion in string-to-uint resolvers).
- Consider replacing repeated `ctx.Value(ctxkeys.UserID)` extraction with a shared helper to reduce duplication.
- Review S2S role bypass in `@auth` and document explicit security constraints for trusted callers.
- Add a short “schema change workflow” section (edit module -> run `make graphql` -> validate resolvers) in this README.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../../README.md)
<!-- doc-nav:end -->
