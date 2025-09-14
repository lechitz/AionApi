# HTTP Utils (Shared)

**Folder:** `internal/platform/server/http/helpers/httpresponse`

## Responsibility

* Provide tiny, reusable helpers for **HTTP cookie** handling.
* Standardize how **auth/session cookies** are created, configured, and cleared across adapters.
* Keep security flags and lifetimes **centralized** (driven by `internal/platform/config`).

## How it works

* Helpers read cookie-related settings from **platform `config`** (name, domain, path, TTL, SameSite, Secure/HTTPOnly).
* Values like the canonical cookie name reuse **shared constants** (e.g., `commonkeys.AuthTokenCookieName`) to avoid drift.
* Utilities write the `Set-Cookie` header using sane defaults for production (HttpOnly + Secure, explicit `Max-Age`/`Expires`, `SameSite`).

## Key helpers

* `SetAuthCookie(w http.ResponseWriter, token string, cfg config.<...>)`

    * Writes a secure, HTTP-only cookie containing the auth token.
    * Applies TTL from config (sets both `Max-Age` and `Expires`).
    * Honors domain/path and SameSite/Secure flags from config.
* (If/when added) **Clear/Unset helper**

    * Writes a past-date cookie with `Max-Age=0` to remove it (useful for logout flows).

> Exact function signatures may evolve; the behavior above is what adapters should rely on.

## Conventions

* **Never** log raw cookie values or tokens.
* Always set:

    * `HttpOnly=true`
    * `Secure=true` (except in strictly local dev over HTTP)
    * An explicit `SameSite` policy (e.g., `Lax` for browser flows; `None` only when you also set `Secure`)
* Prefer **short-lived** cookies and renew on sensitive state changes (e.g., password update).
* Keep cookie names/keys in **`internal/shared/constants`** to ensure consistency.

## Example (illustrative)

```go
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    // ... authenticate user, get signed token ...
    // httputils.SetAuthCookie(w, token, h.Config) // apply platform-wide flags
    // httpresponse.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
```

## Testing hints

* Use `httptest.NewRecorder()` and assert on `rec.Header().Get("Set-Cookie")`:

    * Contains **cookie name** and **value**.
    * Includes `HttpOnly`, `Secure`, `SameSite=<expected>`.
    * Has **`Max-Age`** and **`Expires`** consistent with configured TTL.
* For logout/clear, assert a cookie with `Max-Age=0` and an **expired `Expires`**.

## Design notes

* All cookie logic lives here to keep **handlers thin** and reduce copy-paste of security flags.
* Configuration-driven behavior lets you switch environments (dev/prod) without touching adapters.
