// Package config provides centralized configuration structures and loading logic.
package config

import "time"

// ObservabilityConfig holds all observability-related configuration.
type ObservabilityConfig struct {
	OtelExporterOTLPEndpoint string `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT" default:"otel-collector:4318"`
	OtelServiceName          string `envconfig:"OTEL_SERVICE_NAME"           default:"AionApi"`
	OtelServiceVersion       string `envconfig:"OTEL_SERVICE_VERSION"        default:"0.1.0"`
}

// DBConfig holds database connection configuration.
type DBConfig struct {
	Name     string `envconfig:"DB_NAME"     required:"true"`
	Host     string `envconfig:"DB_HOST"                     default:"localhost"`
	Port     string `envconfig:"DB_PORT"                     default:"5432"`
	Type     string `envconfig:"DB_TYPE"                     default:"postgres"`
	User     string `envconfig:"DB_USER"     required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
}

// CacheConfig holds Redis cache configuration.
type CacheConfig struct {
	Addr     string `envconfig:"CACHE_ADDR"      default:"redis-aion:6379"`
	Password string `envconfig:"CACHE_PASSWORD"`
	DB       int    `envconfig:"CACHE_REPO"      default:"0"`
	PoolSize int    `envconfig:"CACHE_POOL_SIZE" default:"10"`
}

// Secret holds encryption/authentication secrets.
type Secret struct {
	Key string `envconfig:"SECRET_KEY"`
}

// ServerGraphql holds GraphQL server configuration.
type ServerGraphql struct {
	Port string `envconfig:"GRAPHQL_PORT" default:"8081" required:"true"`
}

// ServerHTTP holds HTTP server configuration.
type ServerHTTP struct {
	Port         string        `envconfig:"HTTP_PORT"          default:"5001"     required:"true"`
	Context      string        `envconfig:"HTTP_CONTEXT"       default:"aion-api"`
	ReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT"  default:"10s"`
	WriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"10s"`
}

// Application holds general application-related configuration.
type Application struct {
	ContextRequest time.Duration `envconfig:"CONTEXT_REQUEST"  default:"2.1s"`
	Timeout        int           `envconfig:"SHUTDOWN_TIMEOUT" default:"5"`
}
