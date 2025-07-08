// Package output internal/core/ports/output/cache_output.go
package output

import (
	"context"
	"errors"
	"time"
)

// ErrNil is an error returned when a cache key does not exist.
var ErrNil = errors.New("cache: key does not exist")

// Cache is an abstraction for a cache service.
type Cache interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	Close() error
}
