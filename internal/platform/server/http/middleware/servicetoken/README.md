# HTTP Middleware — Service Token (Platform)

**Folder:** `internal/platform/server/http/middleware/servicetoken`

## Purpose and Main Capabilities

- Authenticate trusted service-to-service calls via `X-Service-Key`.
- Optionally impersonate a user via `X-Service-User-Id`.
- Mark requests as service accounts in context.

## How it works

- If `X-Service-Key` is absent, the request passes through.
- If present, it is validated against `cfg.AionChat.ServiceKey`.
- On success, it sets:
  - `ctxkeys.ServiceAccount = true`
  - `ctxkeys.UserID` if `X-Service-User-Id` is valid.
- On failure, responds with `401 Unauthorized`.

## Headers

- `X-Service-Key` (required for S2S)
- `X-Service-User-Id` (optional)

## Configuration

```env
AION_CHAT_SERVICE_KEY=your-secret-key-here
```

## Usage

Apply to the GraphQL route (or any protected S2S route):

```go
r.GroupWith(servicetoken.New(cfg, log), func(sr ports.Router) {
  sr.Mount(cfg.ServerGraphql.Path, graphqlHandler)
})
```

## Context Keys

- `ctxkeys.ServiceAccount` (`bool`)
- `ctxkeys.UserID` (`uint64`, when provided)

## What Should NOT Live Here

- Domain logic or policy decisions.

