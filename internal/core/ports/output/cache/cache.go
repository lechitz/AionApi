package cache

import "time"

type ICache interface {
	Set(key string, value any, ttl time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Close() error
}
