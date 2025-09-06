// Package cache provides methods for managing tokens in the cache.
package cache

import (
	"github.com/lechitz/AionApi/internal/platform/ports/output/cache"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// Store is a repository for managing tokens.
type Store struct {
	cache  cache.Cache
	logger logger.ContextLogger
}

// New creates a new instance of Store with a given cache and logger.
func New(cache cache.Cache, logger logger.ContextLogger) *Store {
	return &Store{
		cache:  cache,
		logger: logger,
	}
}
