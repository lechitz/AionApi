# internal/chat

Chat domain for message flows exposed via HTTP and backed by external chat/AI integrations.

## Purpose and Main Capabilities

- Accept and process chat messages via HTTP.
- Orchestrate external chat providers through output ports.
- Validate message/session invariants in core.
- Normalize external errors into semantic failures for adapters.

## Package Composition

- `core/`: chat entities, ports, and usecases.
- `core/ports/input`: chat service interface for HTTP handlers.
- `core/ports/output`: contracts for external chat providers and integrations.
- `core/usecase`: send/list/acknowledge message workflows.
- `adapter/primary/http`: handlers, DTO mapping, and middleware usage.
- `adapter/secondary`: integrations with external chat/graph services.

## Flow (Where it comes from -> Where it goes)

HTTP request -> primary adapter -> input port -> usecase ->
output port -> secondary adapter -> external service -> response

## How It Works (Concise)

- Handlers validate input, open spans, and call the chat service (input port).
- Usecases enforce limits, sanitization, and retention policy where applicable.
- Secondary adapters handle timeouts/retries and translate provider errors.

## Why It Was Designed This Way

- Keep core chat rules testable and provider-agnostic.
- Swap external services without changing core logic.
- Centralize reliability and observability at IO boundaries.

## Recommended Practices Visible Here

- Do not log message content; log metadata only.
- Use semantic errors for external failures.
- Keep handlers thin: decode, validate, call usecase, map response.

## Differentials

- Provider-agnostic chat orchestration via output ports.

## What Should NOT Live Here

- Business logic in adapters.
- Provider SDK types in core.
- Cross-context imports.
