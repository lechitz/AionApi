# Fx Application Wiring

**Path:** `internal/platform/fxapp`

## Purpose

This package is the Uber Fx composition root for `aion-api`.
It wires infrastructure providers, application services, HTTP runtime, realtime consumption, and the dedicated outbox publisher process.

## Current Modules

| Module | Role |
| --- | --- |
| `InfraModule` | logger, config, cache, DB, HTTP client, observability init |
| `ApplicationModule` | compose repositories, use cases, and `app.Dependencies` |
| `ServerModule` | compose HTTP handler, build server, and manage lifecycle |
| `RealtimeModule` | start Kafka projection consumer when realtime is enabled |
| `OutboxPublisherModule` | start the periodic Kafka outbox publisher loop |

## Runtime Use

- `cmd/api` boots `InfraModule`, `ApplicationModule`, `RealtimeModule`, and `ServerModule`
- `cmd/outbox-publisher` boots `InfraModule` and `OutboxPublisherModule`

## Boundaries

- Fx wiring belongs here, not in the bounded contexts
- providers should expose stable contracts and delegate behavior to owning packages
- if startup behavior changes, update the matching command README and runtime docs in the same PR

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
