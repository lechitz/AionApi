# internal/adapter/secondary/cache

Cache adapters implementing domain/platform cache ports.

## Package Composition

- Provider-specific cache implementations (e.g., Redis).

## Flow (Where it comes from -> Where it goes)

Usecase -> cache port -> cache adapter -> cache backend

## Why It Was Designed This Way

- Keep caching concerns out of core logic.
- Allow provider swap without changing usecases.

## Recommended Practices Visible Here

- Define key formats and TTLs per domain.
- Handle serialization/deserialization defensively.
- Emit metrics/logs for hit/miss and latency.

## Differentials

- Consistent cache keying patterns.

## What Should NOT Live Here

- Business rules or domain validation.
