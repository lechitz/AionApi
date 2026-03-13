# Event Outbox

This bounded context owns the durable outbox used to persist canonical domain events before they are published to Kafka.

Current scope:
- event enqueue only
- versioned event envelope validation
- DB persistence for later publisher pickup

Out of scope for this package:
- Kafka publishing
- retry scheduler
- projection consumers

The initial integration is intentionally narrow so `AionApi` keeps its role as canonical transactional authority while the event backbone is introduced safely.
