# Cache (Secondary Adapter — Redis)

## Responsibilities

* Implement the cache output port (`internal/platform/ports/output/cache.Cache`) using Redis.
* Hide Redis client details (connection, ping, TTL, deletion) behind a small, testable interface.
* Provide lifecycle control (`Close`) and lightweight health probing (`Ping`).

## How it works

* `NewConnection(appCtx, cfg, logger)` builds a `go-redis` client from `config.CacheConfig` and returns an `output.Cache`.
* Uses a short, bounded context (`cfg.ConnectTimeout`) to `PING` on startup; fails fast if Redis is unreachable.
* Exposes the port methods:

    * `Ping(ctx) error`
    * `Set(ctx, key, value, expiration) error`
    * `Get(ctx, key) (string, error)` — returns `""` on cache miss (maps `redis.Nil` to empty string).
    * `Del(ctx, key) error`
    * `Close() error`
* All methods propagate `context.Context` for deadlines, cancellation, and tracing.

## Reminders

* **Port-only semantics:** never leak `go-redis` types to callers; keep the interface stable.
* **Context-first:** always pass request-scoped `ctx` through to support timeouts and OTel correlation.
* **Cache misses:** treat `redis.Nil` as a benign miss; do not log as an error.
* **Serialization:** `Set` accepts `interface{}`—callers must serialize to a string/JSON as needed; this adapter does not encode structs.
* **Key hygiene:** define key naming/versioning (prefixes, namespaces) in the calling layer; keep the adapter generic.
* **Security:** do not log secrets/values; log minimal metadata (e.g., key prefix, operation, latency if instrumented).
* **TTL discipline:** choose expirations in the use case; avoid “forever” unless intentionally required.
* **Testing:** mock via the `output.Cache` port or use a local Redis in integration tests.
