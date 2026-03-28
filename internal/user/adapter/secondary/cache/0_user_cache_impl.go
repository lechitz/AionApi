// Package cache provides user cache operations.
package cache

import (
	"github.com/lechitz/aion-api/internal/platform/ports/output/cache"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
)

// Store is a repository for managing user profile data in cache.
// SECURITY: This store NEVER caches sensitive data like password hashes.
// Password verification always hits the database.
type Store struct {
	cache  cache.Cache
	logger logger.ContextLogger
}

// NewStore creates a new instance of Store with a given cache and logger.
func NewStore(cache cache.Cache, logger logger.ContextLogger) *Store {
	return &Store{
		cache:  cache,
		logger: logger,
	}
}
