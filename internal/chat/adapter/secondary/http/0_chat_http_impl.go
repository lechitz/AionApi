// Package http provides the HTTP client adapter for communicating with Aion-Chat service.
package http

import (
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/chat/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// AionChatClient is the HTTP client for communicating with Aion-Chat service.
type AionChatClient struct {
	httpClient *http.Client
	baseURL    string
	logger     logger.ContextLogger
}

// NewClient creates a new AionChatClient instance.
func NewClient(baseURL string, timeout time.Duration, log logger.ContextLogger) output.AionChatClient {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	return &AionChatClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
		logger:  log,
	}
}
