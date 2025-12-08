# Service Token Middleware

Middleware for service-to-service (S2S) authentication between AionApi and internal services like Aion-Chat.

## Overview

This middleware allows trusted internal services (like Aion-Chat) to make authenticated calls to AionApi's GraphQL endpoint, executing queries on behalf of a specific user.

## How It Works

```
┌─────────────┐                    ┌────────────────┐
│  Aion-Chat  │───────────────────▶│    AionApi     │
│             │  POST /graphql     │                │
│             │  Headers:          │                │
│             │  X-Service-Key     │                │
│             │  X-Service-User-Id │                │
└─────────────┘                    └────────────────┘
```

### Headers

| Header | Required | Description |
|--------|----------|-------------|
| `X-Service-Key` | Yes | Shared secret key |
| `X-Service-User-Id` | No | User ID for impersonation |

### Behavior

| Scenario | Action |
|----------|--------|
| Header absent | Pass through (normal user auth flow) |
| Header valid | Set `ServiceAccount=true` and inject `UserID` |
| Header invalid | Return `401 Unauthorized` |

## Configuration

### AionApi (.env)

```env
AION_CHAT_SERVICE_KEY=your-secret-key-here
```

### Aion-Chat (.env)

```env
AION_API_SERVICE_KEY=your-secret-key-here
AION_API_GRAPHQL_URL=http://localhost:5001/aion/api/v1/graphql
```

### Generate Secure Key

```bash
openssl rand -base64 32
```

## Context Keys

| Key | Type | Description |
|-----|------|-------------|
| `ctxkeys.ServiceAccount` | `bool` | `true` if authenticated via S2S |
| `ctxkeys.UserID` | `uint64` | User ID (if header provided) |

## Usage Example (Python)

```python
headers = {
    "Content-Type": "application/json",
    "X-Service-Key": os.getenv("AION_API_SERVICE_KEY"),
    "X-Service-User-Id": str(user_id),
}
response = await client.post(url, json={"query": query}, headers=headers)
```

## Manual Testing

```bash
# Valid request
curl -X POST http://localhost:5001/aion/api/v1/graphql \
  -H "Content-Type: application/json" \
  -H "X-Service-Key: your-key" \
  -H "X-Service-User-Id: 1" \
  -d '{"query": "{ categories { id name } }"}'

# Invalid request (returns 401)
curl -X POST http://localhost:5001/aion/api/v1/graphql \
  -H "X-Service-Key: wrong-key" \
  -d '{"query": "{ categories { id } }"}'
```


