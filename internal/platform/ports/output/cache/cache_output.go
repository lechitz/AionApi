// Package cache output ports for cache service.
package cache

import (
	"context"
	"time"
)

// Cache is an abstraction for a cache service.
type Cache interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	Close() error
}
