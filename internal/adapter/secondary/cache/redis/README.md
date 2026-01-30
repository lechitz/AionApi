# internal/adapter/secondary/cache/redis

Redis implementation of cache output ports.

## Package Composition

- Redis client wiring and cache operations.

## Flow (Where it comes from -> Where it goes)

Usecase -> cache port -> Redis adapter -> Redis

## Why It Was Designed This Way

- Provide a fast, shared cache backend.
- Keep Redis specifics out of the core.

## Recommended Practices Visible Here

- Keep TTLs and key namespaces explicit.
- Use timeouts to prevent blocking calls.
- Avoid storing partial or masked data.

## Differentials

- Standardized Redis keying and TTL strategy.

## What Should NOT Live Here

- Domain logic or transport DTOs.
