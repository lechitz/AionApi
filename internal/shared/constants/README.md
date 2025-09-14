# Shared Constants (Cross-cutting)

**Folder:** `internal/shared/constants`

## Responsibility

* Be the **single source of truth** for string keys used across the codebase (logging fields, headers, JWT claims, context keys, tracing attributes, etc.).
* Eliminate **magic strings/typos**, make refactors safer, and keep cross-cutting concerns consistent.

## What it contains

* `claimskeys/` — JWT claim names (e.g., `UserID`, `Exp`, `TokenValue`).
* `commonkeys/` — Common string keys used in logs, HTTP, and DTOs, grouped by area:

    * `general_commonkeys.go` — app metadata and general HTTP keys (e.g., `APIName`, `AppEnv`, `AppVersion`, `XRequestID`).
    * `user_commonkeys.go` — user-related keys (e.g., `User`, `UserID`, `Username`).
    * `token_commonkeys.go` — token fields (e.g., `Token`, `TokenKey`, `AuthTokenCookieName`).
    * `category_commonkeys.go` — category fields (e.g., `Category`, `CategoryID`, `CategoryName`).
* `ctxkeys/` — **type-safe context keys** (custom type `contextKey`) for values injected into `context.Context` (e.g., `UserID`, `Token`, `RequestID`, `TraceID`, `SpanID`).
* `tracingkeys/` — OpenTelemetry span attribute keys (e.g., `HTTPStatusCodeKey`, `RequestIPKey`, `RequestUserAgentKey`).

## How it’s used

* **Logging**

  ```go
  h.Logger.Infow("category created",
      commonkeys.CategoryID, id,
      commonkeys.UserID, userID,
  )
  ```

* **HTTP headers**

  ```go
  w.Header().Set(commonkeys.XRequestID, reqID)
  ```

* **Context propagation**

  ```go
  ctx = context.WithValue(ctx, ctxkeys.RequestID, reqID)
  rid, _ := r.Context().Value(ctxkeys.RequestID).(string)
  ```

* **Tracing attributes**

  ```go
  span.SetAttributes(
      attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK),
      attribute.String(tracingkeys.RequestIPKey, clientIP),
  )
  ```

* **JWT claims**

  ```go
  claims[claimskeys.UserID] = userID
  claims[claimskeys.Exp] = expUnix
  ```

## Design notes

* **No business logic here** — only constants and the minimal `contextKey` type in `ctxkeys`.
* Prefer **domain-specific keys** from `commonkeys/*` instead of ad-hoc strings in handlers, middleware, and repositories.
* `ctxkeys` are **not** header names; they are internal keys for `context.Context`.
* Add new keys **close to their domain** (e.g., new feature → a focused file under `commonkeys/`).

## Conventions

* Keys are **lowercase snake\_case** for log fields and claims; HTTP header constants use canonical header names (e.g., `XRequestID`).
* Keep names **stable**—changing a key breaks dashboards and log queries.
* Avoid duplicate meanings across files; reuse existing keys where possible.

## Testing hints

* Use constants in gomock expectations to avoid brittle string matches:

  ```go
  logger.EXPECT().Infow(gomock.Any(), commonkeys.UserID, gomock.Any()).AnyTimes()
  ```
* In middleware tests, assert the same **request ID** is present in:

    * `X-Request-ID` response header, and
    * `ctx.Value(ctxkeys.RequestID)` inside downstream handlers.
