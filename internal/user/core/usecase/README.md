# User Usecases (Core)

**Folder:** `internal/user/core/usecase`

## Responsibility

* Implement the **User** domain logic behind the input ports (`internal/user/core/ports/input`).
* Orchestrate validations/policies and call output ports (repo, hasher, token provider, auth store).
* Return **domain entities** and **semantic errors** (no transport/ORM details here).

## How it works

* Dependencies are injected with
  `New(userRepository, authStore, tokenProvider, hasher, logger)`.
* Uses only **interfaces** from `core/ports/input` and `core/ports/output` (DIP).
* Each method accepts `context.Context` and instruments **OTel**:

    * **Tracer:** `aionapi.user.create` (`TracerName`)
    * Spans: `Create`, `GetByID`, `FindByUsername`, `FindByEmail`, `ListAll`, `Update`, `UdpateUserPassword`\*, `SoftDelete`
    * Canonical attributes: `operation`, `user_id`, `username`, `email`, `status`

> \*Small TODO: rename span value `UdpateUserPassword` → `UpdateUserPassword`.

## Use cases

* **Create(ctx, cmd)**
  Normalize input (trim/lower email) → check **uniqueness** (username/email) → `Hasher.Hash(password)` → `UserRepository.Create`.
  Errors: `ErrUsernameInUse`, `ErrEmailInUse`, `ErrorToHashPassword`, `ErrorToCreateUser`.

* **GetByID(ctx, userID)**
  Delegates to repo; traces `user_id`.
  Error label: `ErrorToGetSelf`.

* **GetUserByUsername(ctx, username)**
  Delegates to repo.
  Error label: `ErrorToGetUserByUsername`.

* **ListAll(ctx)**
  Delegates to repo; returns non-deleted users.
  Error label: `ErrorToGetSelf` (db failure).

* **UpdateUser(ctx, userID, cmd)**
  Validates `cmd.HasUpdates()` → builds **partial update map** (only provided fields) + `updated_at` → `repo.Update`.
  Errors: `ErrorNoFieldsToUpdate`, `ErrorToUpdateUser`.

* **UpdatePassword(ctx, userID, old, new)**
  `repo.GetByID` → `Hasher.Compare(storedHash, old)` → `Hasher.Hash(new)` → `repo.Update(password)` →
  `tokenProvider.Generate` → `authStore.Save`.
  Errors: `ErrorToGetSelf`, `ErrorToCompareHashAndPassword`, `ErrorToHashPassword`, `ErrorToUpdatePassword`, `ErrorToCreateToken` (including empty token), `ErrorToSaveToken`.

* **SoftDeleteUser(ctx, userID)**
  `authStore.Delete` → `repo.SoftDelete`.
  If token deletion fails, **stop** and return the error.
  Error label: `ErrorToSoftDeleteUser`.

## Rules of thumb

* **No transport or persistence concerns** here (no DTOs, no ORM types).
* Prefer **semantic errors** from `internal/shared/sharederrors` for validation/conflict.
* **Don’t log PII/secrets** (passwords, hashes, tokens); prefer IDs and counts.
* Normalize input that affects business rules (e.g., lower-case email) **inside** the use case.

## Testing hints

* Unit-test with **gomock** ports (no real DB/cache):

    * `create_test.go`: normalization, uniqueness, hash, repo errors.
    * `get_by_id_test.go`, `get_by_username_test.go`: success/not-found/failures.
    * `list_all_test.go`: empty vs non-empty, db error.
    * `update_user_test.go`: partial updates, no-fields, repo error.
    * `update_password_test.go`: compare/hash/update/token generate/save and all error branches.
    * `soft_delete_user_test.go`: token delete first, short-circuit on failure.

## Related to DTOs

* DTOs live in **HTTP adapter**: `internal/user/adapter/primary/http/dto`.
  Controllers map **DTO → input commands**; the **usecase** layer only sees the commands and returns **domain**.
