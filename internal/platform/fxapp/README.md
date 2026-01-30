# internal/platform/fxapp

Fx application graph wiring (modules and lifecycle hooks) that compose platform + domains.

## Purpose and Main Capabilities

- Expose Fx modules (InfraModule, ApplicationModule, ServerModule) used by `cmd/aion-api`.
- Wire constructors and invocations for db, cache, logger, observability, and servers.
- Manage start/stop hooks with graceful shutdown logging.

## Package Composition

- Module definitions and Fx options.
- Lifecycle hooks for startup and shutdown.

## Flow (Where it comes from -> Where it goes)

cmd/aion-api -> fxapp modules -> platform + adapters -> running server

## Why It Was Designed This Way

- Keep wiring explicit and centralized.
- Separate infra wiring from domain code.
- Make lifecycle behavior observable and testable.

## Recommended Practices Visible Here

- Prefer explicit `fx.Options` over hidden globals.
- Keep module boundaries clear (infra vs application vs server).
- Log failures during start/stop with context.

## What Should NOT Live Here

- Business logic or domain rules.
- Adapter implementations tied to a single context.
