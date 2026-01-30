# cmd/aion-api (Entrypoint)

This package is the application entrypoint. It orchestrates process startup, module wiring, and graceful shutdown. It does not contain business rules.

## Package Composition

- `main.go`
  - Boots the app with `fx.New(...)`.
  - Wires infrastructure, application, and server modules.
  - Controls lifecycle start/stop with context and timeout.
- `swagger.go`
  - Applies runtime Swagger metadata (base path, title, version).
  - Keeps docs configuration close to the entrypoint.
- `README.md`
  - Short onboarding guide for this package.

## Flow (Where it comes from -> Where it goes)

The flow starts at the OS process, loads config, wires modules, starts the server, and shuts down on signals.

![AionAPI Entrypoint Flow](../../docs/diagram/images/cmd-aion-api.svg)

Diagram source: `docs/diagram/cmd-aion-api.sequence.txt`

## Boot Sequence (Short)

1) `main.go` builds the Fx app with:
   - `configureSwagger` (base path, title, version from config).
   - `fxapp.InfraModule` (infra/observability/cache/db).
   - `fxapp.ApplicationModule` (domains, ports, use cases).
   - `fxapp.ServerModule` (HTTP + GraphQL).
2) `app.Start()` brings dependencies and servers up; `app.Done()` blocks until shutdown.
3) Graceful stop with a 10s timeout via `app.Stop(...)`.

## Why It Was Designed This Way

- Keep entrypoint thin and predictable.
- Centralize operational concerns (config, observability, shutdown).
- Keep domain boundaries clean by delegating wiring to platform modules.

## Recommended Practices Visible Here

- Small, single-purpose entrypoint functions.
- Module-based composition via Fx.
- Runtime-configured Swagger metadata.
- Graceful shutdown with context timeouts.

## Differentials

- Runtime Swagger configuration avoids stale docs.
- Entry-point-only observability wiring ensures consistent telemetry.
- Explicit shutdown path keeps services reliable in production.

## Extension Guidelines

- Register new modules in `fxapp`, not in `main`.
- Avoid feature flags in `main`; prefer config-driven switches.

## What Should NOT Live Here

- Domain rules or validation.
- HTTP/GraphQL mapping.
- Repository or external IO logic.
