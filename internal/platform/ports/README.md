# Platform Ports (Output)

**Folder:** `internal/platform/ports`

**Subfolders:** `output/cache`, `output/hasher`, `output/keygen`, `output/logger`

## Responsibility

* Define **technology-agnostic output ports** used across contexts (cross-cutting platform services).
* Decouple the **core/domain** from concrete infra (cache, hashing, key generation, logging).
* Enable **secondary adapters** to plug implementations (Redis, bcrypt/argon2, JWT libs, Zap/Logrus, etc.).

## How it works

* Domain and adapters **depend only on these interfaces**.
* Concrete implementations live in platform/secondary packages (or external modules) and are **injected** where needed.
* Unit tests generate mocks from these ports via `make mocks` and keep the domain fully testable.

---

## Interfaces (overview)

### `output/cache.Cache`

```go
Ping(ctx context.Context) error
Set(ctx context.Context, key string, value any, expiration time.Duration) error
Get(ctx context.Context, key string) (string, error)
Del(ctx context.Context, key string) error
Close() error
```

* Abstracts a key–value cache (e.g., Redis).
* `value` is `any` to allow simple strings or JSON-encoded payloads at the edge.

### `output/hasher.Hasher`

```go
Hash(plain string) (string, error)
Compare(hashed, plain string) error
```

* Password hashing and verification (e.g., bcrypt, argon2).
* **Never** compare secrets directly in the domain; always go through this port.

### `output/keygen.Generator`

```go
Generate() (string, error)
```

* Generates secret keys (e.g., for JWT or crypto).
* Used by config/bootstrap to create a transient secret when `SECRET_KEY` is missing.

### `output/logger.ContextLogger`

```go
// printf-style
Infof(format string, args ...any)
Errorf(format string, args ...any)
Debugf(format string, args ...any)
Warnf(format string, args ...any)

// structured
Infow(msg string, keysAndValues ...any)
Errorw(msg string, keysAndValues ...any)
Debugw(msg string, keysAndValues ...any)
Warnw(msg string, keysAndValues ...any)

// context-aware variants (not shown here) are also mocked and used widely
```

* Structured, leveled logging with optional context-aware variants used throughout the codebase.
* Prefer the **structured** methods for telemetry-friendly logs.

---

## Design notes

* Treat these as **stable contracts**; avoid leaking vendor-specific types here.
* Keep interfaces **small and focused** (interface segregation).
* Output ports = things the **domain calls out to** (hashing, cache, logger, etc.).
* Concrete packages must live outside the domain (e.g., `internal/platform/.../adapter/secondary`).

---

## Usage (example)

```go
// in a use case / service constructor
func NewService(repo Repo, cache cache.Cache, hasher hasher.Hasher, log logger.ContextLogger) *Service {
    return &Service{repo: repo, cache: cache, hasher: hasher, log: log}
}

// in a method
func (s *Service) Register(ctx context.Context, plainPwd string) error {
    hashed, err := s.hasher.Hash(plainPwd)
    if err != nil {
        s.log.Errorw("hash.fail", "error", err)
        return err
    }
    // ...
    return nil
}
```

---

## Testing hints

* Generate mocks from ports with:

  ```
  make mocks
  ```

  Mocks are written to `tests/mocks/` (flat package). **Do not edit** generated files; re-run the target instead.
* In test suites, stub behavior on the mocks and assert interactions:

  ```go
  mockHasher.EXPECT().Hash("secret").Return("hashed", nil)
  mockLogger.EXPECT().Infow(gomock.Any(), gomock.Any()).AnyTimes()
  ```
* Avoid coupling tests to real infrastructure; rely on these ports for **fast, deterministic** unit tests.

---

## Extending

* Adding a new platform capability? Create a new interface under `internal/platform/ports/output/<capability>`.
* Keep the surface minimal; prefer **one interface per concern**.
* Provide a secondary adapter in a separate package and **inject** it at composition time.
