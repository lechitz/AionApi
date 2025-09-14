# Test Setup (helpers for unit tests)

**Folder:** `tests/setup/`

## Responsibility

* Provide **ready-to-use test suites** for services (User, Auth, Token, Category).
* Wire **GoMock** controllers + mocks and return the **system under test (SUT)**.
* Offer shared helpers (e.g., **relaxed logger expectations**, **default test user**).

> ⚠️ These builders are meant for tests only. They hide boilerplate and standardize setup/teardown across packages.

---

## How it works

* Each `*ServiceTest(t *testing.T)` function:

    * Creates a `gomock.Controller`.
    * Instantiates the required **output-port mocks** (from `tests/mocks`).
    * Applies **relaxed logger expectations** via `ExpectLoggerDefaultBehavior(...)`.
    * Builds the concrete service (`usecase.NewService(...)`).
    * Returns a typed `*...TestSuite` containing: `Ctrl`, `Ctx`, `Logger`, the mocks, and the `...Service` SUT.

* You are responsible for teardown:

  ```go
  suite := setup.UserServiceTest(t)
  defer suite.Ctrl.Finish()
  ```

---

## Provided suites & helpers

* `UserServiceTest(t)` → `UserServiceTestSuite`

    * Mocks: `UserRepository`, `AuthStore`, `AuthProvider`, `Hasher`, `ContextLogger`
    * SUT: `*user/usecase.Service`
    * Extras: `setup.DefaultTestUser()` (common `domain.User`)

* `AuthServiceTest(t)` → `AuthServiceTestSuite`

    * Mocks: `UserRepository`, `AuthStore`, `AuthProvider`, `Hasher`, `ContextLogger`
    * SUT: `*auth/usecase.Service`

* `TokenServiceTest(t)` → `TokenServiceTestSuite`

    * Mocks: `AuthStore`, `AuthProvider`, `ContextLogger` (and internal deps wired for the SUT)
    * SUT: `*auth/usecase.Service`

* `CategoryServiceTest(t)` → `CategoryServiceTestSuite`

    * Mocks: `CategoryRepository`, `ContextLogger`
    * SUT: `*category/usecase.Service`

* `ExpectLoggerDefaultBehavior(logger *mocks.MockContextLogger)`

    * Registers **non-intrusive** expectations (`DebugwCtx/InfowCtx/WarnwCtx/ErrorwCtx`) with multiple arities so tests don’t fail on logging KVs.

---

## Usage examples

```go
func TestGetUserByID_Success(t *testing.T) {
    suite := setup.UserServiceTest(t)
    defer suite.Ctrl.Finish()

    u := setup.DefaultTestUser()
    suite.UserRepository.
        EXPECT().
        GetByID(gomock.Any(), u.ID).
        Return(u, nil)

    got, err := suite.UserService.GetByID(suite.Ctx, u.ID)
    require.NoError(t, err)
    require.Equal(t, u, got)
}
```

```go
func TestSoftDeleteUser_TokenError(t *testing.T) {
    suite := setup.UserServiceTest(t)
    defer suite.Ctrl.Finish()

    suite.TokenStore.
        EXPECT().
        Delete(gomock.Any(), uint64(1)).
        Return(errors.New("delete token failed"))

    err := suite.UserService.SoftDeleteUser(suite.Ctx, 1)
    require.Error(t, err)
}
```

---

## Conventions & tips

* **Always** `defer suite.Ctrl.Finish()` to assert mock expectations.
* Keep assertions focused; use `gomock.AssignableToTypeOf(...)` for dynamic maps/structs.
* Use `tests/mocks` **generated** doubles—don’t hand-edit mocks.
* Prefer **domain-centric** fixtures (`DefaultTestUser`) and keep test data minimal.

---

## When to extend

* Add a new service? Provide a `XServiceTest(t)` builder returning a `XServiceTestSuite`.
* Need new shared fixtures or expectations? Add helpers here to keep tests consistent across contexts.
