// Package httpclient provides an abstraction for HTTP client operations.
package httpclient

import (
	"context"
	"net/http"
)

// HTTPClient is an abstraction for performing HTTP requests with instrumentation support.
// This interface allows adapters to depend on an abstraction rather than a concrete *http.Client,
// following the Dependency Inversion Principle and Hexagonal Architecture patterns.
//
// Implementations should provide:
// - Automatic tracing and span creation (OTEL instrumentation)
// - Context propagation (trace headers injection)
// - Configurable timeouts and transport options.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(ctx context.Context, url string) (*http.Response, error)
	Post(ctx context.Context, url, contentType string, body interface{}) (*http.Response, error)
}
