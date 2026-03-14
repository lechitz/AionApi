# Chat Bounded Context

**Path:** `internal/chat`

## Overview

Chat domain integrating conversational workflows with user data context.
Orchestrates chat-related usecases and external AI/chat providers via output ports.

## Typical Responsibilities

| Area | Responsibility |
| --- | --- |
| Chat history/context | Query and format user conversation context |
| External provider orchestration | Call AI/chat service adapters via output ports |
| Transport integration | Expose chat usecases through primary adapters |
| Audit publishing | Emit audit action events through `audit` bounded context |

## Design Notes

- Keep provider-specific details in secondary adapters.
- Keep message-handling rules in core usecases.
- Avoid logging raw sensitive content.
- Keep audit persistence out of `chat` DB adapters; use `audit` input port/service.

## Package Improvements

- Add provider fallback/retry behavior documentation.
- Add test scenarios for timeout and degraded provider responses.
- Add data retention guidance for chat history artifacts.
- Add per-operation observability key recommendations.

---

<!-- doc-nav:start -->
## Navigation
- [Back to parent layer](../README.md)
- [Back to root README](../../README.md)
<!-- doc-nav:end -->
