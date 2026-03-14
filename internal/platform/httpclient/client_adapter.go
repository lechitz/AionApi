// Package httpclient provides an instrumented HTTP client implementation.
package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	output "github.com/lechitz/AionApi/internal/platform/ports/output/httpclient"
)

// instrumentedClient wraps *http.Client and implements the output.HTTPClient interface.
type instrumentedClient struct {
	client *http.Client
}

// NewClient wraps an *http.Client and returns it as the output.HTTPClient interface.
// The provided client should be instrumented with OTEL (e.g., using otelhttp.NewTransport).
func NewClient(client *http.Client) output.HTTPClient {
	if client == nil {
		// Fallback to default client if none provided
		client = http.DefaultClient
	}
	return &instrumentedClient{client: client}
}

// Do send an HTTP request and returns an HTTP response.
func (c *instrumentedClient) Do(req *http.Request) (*http.Response, error) {
	// #nosec G107,G704 -- The caller is responsible for constructing and validating the outbound request target.
	return c.client.Do(req)
}

// Get issues a GET request to the specified URL.
func (c *instrumentedClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf(errMsgCreateGETRequest, err)
	}
	return c.client.Do(req)
}

// Post issues a POST request to the specified URL with the given content type and body.
func (c *instrumentedClient) Post(ctx context.Context, url, contentType string, body interface{}) (*http.Response, error) {
	var bodyReader *bytes.Buffer
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf(errMsgMarshalBody, err)
		}
		bodyReader = bytes.NewBuffer(jsonData)
	} else {
		bodyReader = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf(errMsgCreatePOSTRequest, err)
	}
	req.Header.Set(headerContentType, contentType)

	return c.client.Do(req)
}
