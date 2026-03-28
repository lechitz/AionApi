# Platform Configuration

**Path:** `internal/platform/config`

## Purpose

This package is the canonical source for typed runtime configuration, env defaults, and cross-section validation used to bootstrap `aion-api`.

## Canonical Source

- field definitions, env tags, and defaults live in `environments.go`
- cross-field validation rules live in `config.go`

When docs conflict with those files, the code wins.

## Current Config Sections

| Section | What it controls |
| --- | --- |
| `General` | app name, environment, version |
| `Observability` | OTLP endpoint, service identity, exporter headers/compression/timeouts |
| `ServerHTTP` | HTTP host/port, context, API root, Swagger/docs/health paths, timeouts |
| `ServerGraphql` | GraphQL host/path and transport limits |
| `DB` | PostgreSQL connectivity and pool/retry settings |
| `Cache` | Redis address, DB isolation by bounded context, pool/timeout |
| `Kafka` | broker list and canonical topic names |
| `Outbox` | batch size, publish interval, enabled flag |
| `Realtime` | SSE path, consumer-group prefix, heartbeat, subscriber buffer |
| `Cookie` | auth cookie domain/path/samesite/secure/max-age |
| `AionChat` | external `aion-chat` base URL, service key, timeout |
| `AvatarStorage` | S3-compatible avatar storage configuration |
| `Application` | shutdown timeout and request context timeout |

## Validation Coverage

`Config.Validate()` currently enforces:

- HTTP and GraphQL path/timeout/header constraints
- cache and DB minimums/required fields
- observability endpoint and compression rules
- Kafka topic/broker requirements
- outbox and realtime runtime minimums
- application shutdown/runtime constraints

## Rules

- Add new env keys by extending the typed structs first, then validation if needed.
- Do not duplicate long env key tables in distant READMEs; link back here instead.
- Any contract-visible change to paths, topics, or external endpoints must update the nearest consumer/operator docs in the same PR.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../../README.md)
<!-- doc-nav:end -->
