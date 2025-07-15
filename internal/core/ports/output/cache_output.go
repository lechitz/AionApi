// Package output internal/core/ports/output/cache_output.go
package output

import (
	"context"
	"time"
)

// Cache is an abstraction for a cache service.
type Cache interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, tokenKey string) (string, error)
	Del(ctx context.Context, key string) error
	Close() error
}
