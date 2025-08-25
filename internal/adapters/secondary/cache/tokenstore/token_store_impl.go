// Package tokenstore provides methods for managing tokens in the cache.
package tokenstore

import (
	"github.com/lechitz/AionApi/internal/core/ports/output"
)

// Store is a repository for managing tokens.
type Store struct {
	cache  output.Cache
	logger output.ContextLogger
}

// New creates a new instance of Store with a given cache and logger.
func New(cache output.Cache, logger output.ContextLogger) *Store {
	return &Store{
		cache:  cache,
		logger: logger,
	}
}
