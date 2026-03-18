# Auth Bounded Context

**Path:** `internal/auth`

## Purpose

`internal/auth` owns login, session validation, refresh-token rotation, logout, and role/session cache integration.

It is the backend authority for authenticated HTTP session state used by REST handlers and GraphQL auth enforcement.

## Current HTTP Surface

| Route | Access | Current behavior |
| --- | --- | --- |
| `POST /auth/login` | public | validates credentials, returns access token + user snapshot, sets auth and refresh cookies |
| `POST /auth/refresh` | public | rotates access and refresh tokens from the refresh cookie |
| `GET /auth/session` | authenticated | validates bearer token or auth cookie and returns session snapshot |
| `POST /auth/logout` | authenticated | revokes current server-side session state and clears cookies |

## Internal Shape

| Area | Responsibility |
| --- | --- |
| `core/usecase` | `Login`, `Validate`, `RefreshTokenRenewal`, `Logout`, and role-cache orchestration |
| `core/ports/output` | auth provider, auth store, roles reader, role cache contracts |
| `adapter/secondary/cache` | Redis-backed token/session and role-cache operations |
| `adapter/primary/http/handler` | REST transport mapping for `/auth/*` routes |
| `adapter/primary/http/middleware` | protected-route middleware that injects authenticated user context |

## Boundaries

- Cookie transport rules are owned by `internal/platform/server/http/utils/cookies`, not by core usecases.
- GraphQL `@auth` enforcement lives in the central GraphQL adapter, but validation semantics ultimately come from this context.
- Security-sensitive integrations stay behind output ports; transport adapters only decode, map, and emit cookies/responses.

## Related Docs

- [`../platform/server/http/utils/cookies/README.md`](../platform/server/http/utils/cookies/README.md)
- [`../platform/server/http/middleware/servicetoken/README.md`](../platform/server/http/middleware/servicetoken/README.md)
- [`../adapter/primary/graphql/README.md`](../adapter/primary/graphql/README.md)

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
