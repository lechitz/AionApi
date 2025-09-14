# Test Data (shared fixtures for unit tests)

**Folder:** `tests/testdata/`

## Responsibility

* Centralize **reusable, deterministic fixtures** for tests across contexts (user, auth, category).
* Keep test code **small and readable** by avoiding repeated inline structs.
* Provide **safe placeholders** (no real secrets/PII).

---

## What’s inside

Current helpers you can import directly in tests:

* **Auth**

    * `SecretKey` — dummy symmetric key for token-related tests.
    * `TestPerfectToken` — a ready-to-use `auth.domain.Auth` sample.

* **User**

    * `TestPerfectUser` — fully-populated `user.domain.User` (IDs, timestamps, etc.).
    * `TestPerfectLoginInputUser` — minimal login payload (username/password).
    * `HashedPassword` — placeholder hash string for comparisons/mocks.

* **Category**

    * `PerfectCategory` — representative `category.domain.Category` sample.

> These are **examples** intended to be copied/adjusted in your tests as needed. They should never reflect production data.

---

## How to use (examples)

```go
import (
    "testing"

    "github.com/lechitz/AionApi/tests/testdata"
    "github.com/stretchr/testify/require"
)

func TestSomething(t *testing.T) {
    u := testdata.TestPerfectUser
    require.NotZero(t, u.ID)
}

func TestLoginPayload(t *testing.T) {
    in := testdata.TestPerfectLoginInputUser
    require.NotEmpty(t, in.Username)
    require.NotEmpty(t, in.Password)
}
```

---

## Conventions & guidelines

* **Determinism first.** Prefer fixed timestamps and IDs in fixtures so your tests don’t depend on `time.Now()`.
  If you need “now”, set it **in the test** (not in shared fixtures) to keep control and predictability.

* **Domain only.** Fixtures should use **core domain types** (e.g., `internal/<context>/core/domain`).
  Avoid importing adapters/transport packages.

* **No real secrets or PII.** Keep keys, tokens, emails, and names obviously fake.

* **Small & composable.** Provide minimal-but-realistic objects. Let each test adjust the fields it actually cares about.

* **Don’t overfit.** If a test requires a very specific shape, build the struct inline in the test or add a tiny helper inside that test file.

---

## Adding new fixtures

1. Create a new file under `tests/testdata/` named by context (e.g., `project.go`).
2. Use clear names: `TestPerfect<Project>`, `Minimal<Project>`, `WithInvalid<Field>`, etc.
3. Keep values **static** (not derived at runtime) unless a test must set them dynamically.
4. If you introduce sensitive-like fields (tokens/passwords), make them obviously dummy.

---

## Notes

* Keep imports aligned with the current domain packages (if package paths move, update fixtures accordingly).
* Unlike `tests/mocks/`, these files are **not generated**—they’re meant to be edited and reviewed.
* If a fixture starts to introduce coupling or flakiness, prefer **local fixtures** in the test file instead.
