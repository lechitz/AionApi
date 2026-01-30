# Platform HTTP

**Folder:** `internal/platform/server/http`

## Responsibility

* Compose the HTTP server, routing, and middleware used by all contexts.
* Provide framework-agnostic router ports and concrete router adapters.
* Expose generic handlers (health, 404/405, error) and shared utilities.

## Package Composition

- `ports/`: router contract used by primary adapters.
- `router/chi/`: chi implementation of the router port.
- `middleware/`: request ID, recovery, cors, service-token.
- `generic/`: health, not-found, error handlers.
- `utils/`: httpresponse, sharederrors, cookies helpers.
- `errors/`: shared HTTP error definitions.
- `composer.go`: mounts platform endpoints and context routes.
- `server.go`: builds and runs the HTTP server.

## Flow (Where it comes from -> Where it goes)

HTTP request -> router -> middleware -> handler -> response helpers

## How It Works (Concise)

- The composer picks the concrete router adapter and applies middlewares in order.
- Context adapters register their routes through `ports.Router` only.
- Generic handlers standardize health and error responses.
- Utilities centralize response envelopes, error mapping, and cookies.

## Recommended Practices Visible Here

- Apply `requestid` early; keep `recovery` outermost.
- Keep handlers thin and transport-only.
- Map semantic errors to HTTP consistently via utils.
- Never import `chi` directly in contexts.

## What Should NOT Live Here

- Domain rules or usecase orchestration.
- Context-specific DTOs or handlers.
