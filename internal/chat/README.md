# Chat Bounded Context

**Path:** `internal/chat`

## Purpose

`internal/chat` owns authenticated chat interaction flows, conversation history/context retrieval, and the integration boundary with the external `aion-chat` service.

## Current Transport Surface

| Surface | Current contract |
| --- | --- |
| HTTP `POST /chat/text` | authenticated message send; returns assistant response, UI payload, sources, and optional usage |
| HTTP `POST /chat/cancel` | authenticated cancel request; proxies cancellation to `AION_CHAT_URL/internal/cancel` |
| HTTP `POST /chat/audio` | authenticated voice/chat ingestion path |
| GraphQL read surface | context controllers expose chat history and aggregated chat context |

There is no shared GraphQL mutation contract for chat in the current backend-owned public surface.

## Runtime Flow

1. Primary adapters authenticate the request and extract `userID`.
2. `ProcessMessage` loads up to six recent cached chat messages for context.
3. The secondary HTTP adapter forwards the request to `aion-chat`.
4. The usecase maps response text, UI payload, sources, and token usage into the local domain result.
5. UI-action audit persistence is attempted best-effort through the `audit` bounded context.
6. Chat history is saved asynchronously with `context.WithoutCancel` so request cancellation does not drop persistence work.

## Internal Shape

| Area | Responsibility |
| --- | --- |
| `core/usecase` | `ProcessMessage`, `GetChatHistory`, `GetChatContext`, async history save, audit event emission |
| `adapter/secondary/http` | external `aion-chat` client boundary |
| `adapter/secondary/db` | durable chat history persistence |
| `adapter/secondary/cache` | cached conversation history access |
| `adapter/primary/http` | authenticated write/cancel/audio endpoints |
| `adapter/primary/graphql/controller` | read-only GraphQL controller surface |

## Boundaries

- Provider-specific HTTP semantics stay in the secondary adapter.
- Audit persistence failures must not fail the main chat response path.
- UI-action metadata belongs to request context/contracts; business ownership of audit storage remains in `internal/audit`.

## Related Docs

- [`../audit/README.md`](../audit/README.md)
- [`../adapter/primary/graphql/README.md`](../adapter/primary/graphql/README.md)
- Cross-repo reference: `aion-docs/planning/v1/reference/aion-chat-current-state.md`

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
