// Package cache provides methods for managing chat history in the cache.
package cache

import (
	"fmt"

	output "github.com/lechitz/aion-api/internal/platform/ports/output/cache"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
)

// Store is a repository for managing chat history in cache.
type Store struct {
	cache  output.Cache
	logger logger.ContextLogger
}

// NewStore creates a new instance of Store with a given cache and logger.
func NewStore(cache output.Cache, logger logger.ContextLogger) *Store {
	return &Store{
		cache:  cache,
		logger: logger,
	}
}

// buildKey builds a cache key for a user's chat history.
func buildKey(userID uint64) string {
	return fmt.Sprintf("%s%d", cacheKeyPrefix, userID)
}
