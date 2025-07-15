// Package output token defines interfaces for token-related operations.
package output

import (
	"context"
	"net/http"
)

// TokenClaimsExtractor defines an interface for extracting claims from a token (JWT, etc.)
type TokenClaimsExtractor interface {
	ExtractFromRequest(r *http.Request) (map[string]interface{}, error)
	ExtractFromContext(ctx context.Context) (map[string]interface{}, error)
}
