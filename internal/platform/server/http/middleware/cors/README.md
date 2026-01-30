# HTTP Middleware — CORS (Platform)

**Folder:** `internal/platform/server/http/middleware/cors`

## Purpose and Main Capabilities

- Apply CORS rules for the frontend origins used by Aion.
- Allow credentialed requests with explicit origins.
- Expose `Set-Cookie` to support auth/session cookies.

## How it works

- `New()` returns a `func(http.Handler) http.Handler` using `github.com/go-chi/cors`.
- Configured with:
  - AllowedOrigins: `http://localhost:5000`, `http://localhost:5173`
  - AllowedMethods: `GET, POST, PUT, DELETE, OPTIONS`
  - AllowedHeaders: `Accept, Authorization, Content-Type, X-CSRF-Token`
  - ExposedHeaders: `Set-Cookie`
  - AllowCredentials: `true`
  - MaxAge: `300`

## Usage

Register in the HTTP composer:

```go
// composer.go
r.Use(cors.New())
```

## Notes

- Keep origins aligned with frontend environments.
- Avoid wildcard origins when `AllowCredentials` is true.

## What Should NOT Live Here

- Domain logic or auth rules.
