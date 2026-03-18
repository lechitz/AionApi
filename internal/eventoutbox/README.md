# Event Outbox

## Purpose

`internal/eventoutbox` owns the durable outbox that persists canonical backend events before publication to Kafka.

It is the relay boundary between transactional writes inside `AionApi` and the wider event backbone.

## Current Surface

| Interface | Responsibility |
| --- | --- |
| `core/ports/input.Service.Enqueue` | validate and persist one outbox event |
| `core/ports/input.PublisherService.PublishPending` | publish one batch of pending events and mark or reschedule rows |
| `core/ports/output.EventRepository` | save, list pending, mark published, reschedule, and expose stats |
| `adapter/secondary/kafka` | publish normalized outbox envelopes to Kafka |

## Runtime Contract

- durable rows are stored in `aion_api.event_outbox`
- newly enqueued events use the backend-owned canonical envelope/version defaults
- the publisher loop reads pending rows in batches, publishes externally, then either:
  - marks rows as published, or
  - reschedules them with backoff and last-error metadata
- aggregate stats are available through repository support code for operator diagnostics

## Boundaries

- Producer contexts own business semantics and decide when an event should be enqueued.
- `eventoutbox` owns durability and publication mechanics, not business behavior.
- Consumers, projections, realtime fanout, and downstream retries are outside this bounded context.
- This package is not directly exposed through REST or GraphQL.

## Related Docs

- Cross-repo reference: `aion-docs/planning/v1/reference/event-backbone-baseline.md`
- [`../platform/config/README.md`](../platform/config/README.md)
