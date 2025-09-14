# Auth Use Cases (Core)

## Responsibilities

* Implement the authentication business flows for the **auth** context.
* Provide a single service that:

    * **Login**: verifies credentials, **issues a token**, and stores its reference.
    * **Logout**: **revokes** the stored token reference.
    * **Validate**: verifies token integrity, extracts `userID`, and checks **cache consistency** (mismatch protection).

## What it provides (operations)

* `Login(ctx, username, password) (user, token string, err error)`
* `Logout(ctx, userID uint64) error`
* `Validate(ctx, rawToken string) (userID uint64, claims map[string]any, err error)`

## Dependencies (ports)

* `userOutput.UserRepository` — look up users (e.g., `GetByUsername`).
* `authOutput.AuthProvider` — **issue/verify** tokens (e.g., JWT).
* `authOutput.AuthStore` — **persist/retrieve/delete** token references (cache).
* `hasher.Hasher` — password comparison.
* `logger.ContextLogger` — structured logging with context.

*All dependencies are **interfaces**; concrete adapters are injected by `internal/platform/bootstrap` (DIP).*

## How it works (high-level)

### Login

1. Lookup user by username.
2. Compare stored hash vs. provided password.
3. Generate token via `AuthProvider`.
4. Save `{Key: userID, Token: <token>}` to `AuthStore`.
5. Return `(user, token)`.

**Typical errors**: `error to get user by username`, `invalid credentials`, `error to create token`.

### Logout

1. Delete token reference from `AuthStore` by `userID`.
2. Return `nil` on success.

**Typical errors**: `error to delete token`.

### Validate

1. Sanitize header value (`Bearer <token>` → `<token>`).
2. Verify signature/exp via `AuthProvider` → `claims`.
3. Extract `userID` from claims (`userID` or fallback `sub`).
4. Get stored token from `AuthStore` and **compare** with provided value.
5. Return `(userID, claims)` on success.

**Typical errors**: `invalid access reference`, `invalid userID in claims`,
`error retrieving access reference from cache`, `provided reference does not match stored one`.

## Observability

* Tracer name: `aionapi.auth`.
* Spans: `Login`, `Logout`, `ValidateToken`.
* Emits events (e.g., `lookup_user`, `compare_password`, `generate_token`, `save_token_to_store`, `verify_token`, `compare_token`).
* Logs use request context (request ID, trace/span IDs, user ID when available). **Never** log token values.

## Notes & Reminders

* **Stateless & concurrent**: the service keeps no mutable state; safe for concurrent use.
* **Do not leak transport/infra**: no HTTP/GraphQL or ORM types here; only domain and ports.
* **Error semantics**: return **meaningful**, domain-level errors; avoid raw strings from infra.
* **Security**: never log secrets (tokens/passwords); log only minimal metadata.
* **Testing**: favor unit tests with mocks of all ports; table-test edge cases (not found, bad password, provider/store failures, mismatch).
