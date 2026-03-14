// Package httpclient provides constants for HTTP client operations.
package httpclient

const (
	// HTTP header keys.
	headerContentType = "Content-Type"

	// Error message formatting.
	errMsgCreateGETRequest  = "failed to create GET request: %w"
	errMsgCreatePOSTRequest = "failed to create POST request: %w"
	errMsgMarshalBody       = "failed to marshal request body: %w"
)
