# API Entrypoint (`cmd/api`)

## Purpose

This package is the process entrypoint for the API binary. It owns bootstrap concerns only.

## File Responsibilities

- `main.go`: process entrypoint only (`main` + `run`).
- `bootstrap_runtime.go`: runtime bootstrap flow (load config, start app, wait signal, graceful stop).
- `bootstrap_fx.go`: Fx composition root (`newFXApp`) with platform modules.
- `bootstrap_config.go`: env-driven bootstrap timeouts (`BOOTSTRAP_START_TIMEOUT`, `BOOTSTRAP_STOP_TIMEOUT`).
- `swagger.go`: Swagger annotation block and runtime Swagger metadata wiring from loaded config.

## When To Change Each File

- `main.go`: only when entrypoint invocation itself changes.
- `bootstrap_runtime.go`: startup/shutdown lifecycle, signal handling, bootstrap logging.
- `bootstrap_fx.go`: dependency graph/modules that compose the application process.
- `bootstrap_config.go`: bootstrap env vars and validation rules.
- `swagger.go`: swagger metadata annotation and runtime swagger metadata mapping.

## Boundaries

- Keep domain logic out of `cmd/api`.
- Keep environment/runtime concerns in `cmd/api`, not in domain packages.
- Keep shared runtime wiring in `internal/platform/fxapp`.
