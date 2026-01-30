# Instrumented HTTP Client

This package provides a helper to create outbound HTTP clients instrumented with OpenTelemetry.

## Overview

The instrumented client automatically:
- Creates spans for each outbound HTTP request
- Injects trace propagation headers (traceparent, tracestate) following W3C Trace Context
- Records HTTP attributes (method, URL, status code) in spans
- Enables distributed tracing across service boundaries

## Architecture & DI

The client follows **Dependency Inversion Principle** (SOLID):
- `ProvideHTTPClient` in `internal/platform/fxapp/infra.go` creates the instrumented client
- Fx provides the client to adapters that need outbound HTTP calls
- Adapters receive `*http.Client` as constructor dependency (not `NewClient(baseURL string, timeout)`)

## Usage

### Via Dependency Injection (Recommended)

Adapters should accept `*http.Client` in constructor:

```go
// Adapter constructor accepts instrumented client
func NewClient(httpClient *http.Client, baseURL string, log logger.ContextLogger) output.Client {
    return &clientImpl{
        httpClient: httpClient,
        baseURL:    baseURL,
        logger:     log,
    }
}

// In Fx wiring (domain.go or infra.go)
func ProvideAppDependencies(httpClient *http.Client, ...) *AppDependencies {
    myClient := adapter.NewClient(httpClient, cfg.Service.BaseURL, log)
    // ...
}
```

### Direct Usage (Tests or Standalone)

```go
client := httpclient.NewInstrumentedClient(httpclient.Options{
    Timeout: 10 * time.Second,
})
resp, err := client.Get("https://example.com")
```

## Configuration

Options:
- `Timeout`: request timeout (default: 15s)
- `BaseTransport`: custom http.RoundTripper (default: http.DefaultTransport)
- `DisableInstrumentation`: disable OTEL for testing
- `OtelOptions`: additional otelhttp.Transport options

## Example: Chat Adapter

The `chat` bounded context uses the instrumented client:

```go
// internal/chat/adapter/secondary/http/0_chat_http_impl.go
func NewClient(httpClient *http.Client, baseURL string, log logger.ContextLogger) output.AionChatClient {
    return &AionChatClient{
        httpClient: httpClient, // Instrumented client from Fx
        baseURL:    baseURL,
        logger:     log,
    }
}
```

Fx wiring:
```go
// internal/platform/fxapp/domain.go
chatHTTPClient := chatClient.NewClient(httpClient, cfg.AionChat.BaseURL, log)
```

## Benefits

1. **Automatic distributed tracing** — spans connect across services
2. **Consistent instrumentation** — all outbound HTTP uses same pattern
3. **SOLID compliance** — adapters depend on abstractions (http.Client interface)
4. **Testability** — inject mock/test clients easily
5. **Configuration centralized** — timeout/transport/headers in one place

## Diagram

![Platform HTTP Client Flow](../../docs/diagram/images/internal-platform-httpclient.svg)

Source: `../../docs/diagram/internal-platform-httpclient.sequence.txt`

## Flow (Where it comes from -> Where it goes)

Adapter -> httpclient (instrumented) -> external HTTP service

## What Should NOT Live Here

- Domain logic or adapter orchestration.
- Service-specific clients (those live in secondary adapters).
