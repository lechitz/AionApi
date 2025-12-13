// Package http provides the HTTP client adapter for communicating with Aion-Chat service.
//
//revive:disable:var-naming // package name intentionally mirrors the protocol we integrate with
package http

//revive:enable:var-naming

import (
	"github.com/lechitz/AionApi/internal/chat/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/ports/output/httpclient"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// AionChatClient is the HTTP client for communicating with Aion-Chat service.
// Follows Hexagonal Architecture: depends on HTTPClient interface (port) instead of concrete implementation.
type AionChatClient struct {
	httpClient httpclient.HTTPClient
	baseURL    string
	logger     logger.ContextLogger
}

// NewClient creates a new AionChatClient instance.
// Accepts an HTTPClient interface as a dependency (Dependency Inversion Principle).
// The caller should provide an instrumented client with appropriate timeout, transport, and OTEL instrumentation.
func NewClient(httpClient httpclient.HTTPClient, baseURL string, log logger.ContextLogger) output.AionChatClient {
	if httpClient == nil {
		panic("http client cannot be nil - ensure platform provider is registered in Fx")
	}

	return &AionChatClient{
		httpClient: httpClient,
		baseURL:    baseURL,
		logger:     log,
	}
}
