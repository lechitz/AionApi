# Bootstrap (Domain Composition Root)

**Path:** `internal/platform/bootstrap`

## Responsibilities

* Initialize **secondary adapters** (DB, cache, token provider, password hasher, etc.).
* Construct and wire **core use cases** (services) that implement **input ports**.
* Expose a single dependency bundle: `*AppDependencies` (e.g., `AuthService`, `UserService`, `CategoryService`, `Logger`) for the Platform layer to mount transports (HTTP/GraphQL).
* Provide a **cleanup** function to gracefully close infrastructure resources.

## How it works

* Entry point:

  ```go
  func InitializeDependencies(
      appCtx context.Context,
      cfg *config.Config,
      logger logger.ContextLogger,
  ) (*AppDependencies, func(context.Context), error)
  ```
* Reads typed configuration from `internal/platform/config`.
* Builds infra:

    * **Cache**: `adapter/secondary/cache.NewConnection(...)`
    * **DB**: `adapter/secondary/db.NewConnection(...)`
    * **Token provider**: `adapter/secondary/token.New(...)`
    * **Hasher**: `adapter/secondary/hasher.New()`
* Wires **repositories** (e.g., `userRepo.NewUser`, `categoryRepo.NewCategory`) and **stores** (e.g., `adapterCache.New`).
* Composes **use cases**:

    * `auth.NewService(userRepo, tokenStore, tokenProvider, hasher, logger)`
    * `user.NewService(userRepo, tokenStore, tokenProvider, hasher, logger)`
    * `category.NewService(categoryRepo, logger)`
* Returns `AppDependencies` + `cleanup(ctx)` that closes DB/Cache respecting the provided context (timeouts/cancel).

## Why this matters

* **Single place** for wiring infra and domain keeps the rest of the app clean and decoupled.
* Enforces **DIP/Hexagonal**: domain depends on ports, not concrete tech.
* Simplifies **testing** and **evolution** (swap adapters without touching handlers/use cases).

## Boundaries

* ❌ Does **not** create routers, handlers, or mount routes.
* ❌ Does **not** read raw env vars directly (only via `config`).
* ✅ Can log structured messages; **never** logs sensitive data.

## Cleanup & lifecycle

* Returns a `cleanup(ctx)` function that:

    * Closes the DB connection (`db.Close`).
    * Closes the cache client (`cacheClient.Close()`).
    * Honors the provided `ctx` for graceful shutdown and logs completion/abort.

## Extending (adding a new dependency)

1. Implement the **output port** and its adapter (e.g., a new repository or external client).
2. Instantiate it here using `cfg` and `logger`.
3. Inject it into the appropriate **use case** constructor.
4. Expose the resulting input port via `AppDependencies`.

## Testing tips

* Unit-test **use cases** with mocks (no need to touch bootstrap).
* For bootstrap itself, use **integration/smoke tests** to ensure connections succeed under test config.
* Verify `cleanup(ctx)` actually releases resources (e.g., no leaked connections).
