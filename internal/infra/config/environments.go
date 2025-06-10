// Package config provides centralized configuration structures and variables.
// for the AionApi application environment setup.
package config

import "time"

// Setting stores the global configuration for the application.
var Setting Config

// Config contains all main configuration sections required to initialize the application.
type Config struct {
	DB            DBConfig
	Cache         CacheConfig
	ServerHTTP    ServerHTTP
	ServerGraphql ServerGraphql
	Application   Application
	Secret        Secret
}

// DBConfig holds the database configuration used by the application.
type DBConfig struct {
	Name     string `envconfig:"DB_NAME"     required:"true"`
	Host     string `envconfig:"DB_HOST"                     default:"localhost"`
	Port     string `envconfig:"DB_PORT"                     default:"5432"`
	Type     string `envconfig:"DB_TYPE"                     default:"postgres"`
	User     string `envconfig:"DB_USER"     required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
}

// CacheConfig contains the cache (Redis) configuration for the application.
type CacheConfig struct {
	Addr     string `envconfig:"CACHE_ADDR"      default:"redis-aion:6379"`
	Password string `envconfig:"CACHE_PASSWORD"`
	DB       int    `envconfig:"CACHE_REPO"      default:"0"`
	PoolSize int    `envconfig:"CACHE_POOL_SIZE" default:"10"`
}

// ServerGraphql holds configuration settings for the GraphQL server.
type ServerGraphql struct {
	Port string `envconfig:"GRAPHQL_PORT" default:"8081" required:"true"`
}

// ServerHTTP holds configuration settings for the HTTP server.
type ServerHTTP struct {
	ReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT"  default:"10s"`
	WriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"10s"`
	Context      string        `envconfig:"HTTP_CONTEXT"       default:"aion-api"`
	Port         string        `envconfig:"HTTP_PORT"          default:"5001"     required:"true"`
}

// Application holds general application-related configuration.
type Application struct {
	Timeout        int           `envconfig:"SHUTDOWN_TIMEOUT" default:"5"`
	ContextRequest time.Duration `envconfig:"CONTEXT_REQUEST"  default:"2.1s"`
}

// Secret holds secret keys used for encryption or authentication.
type Secret struct {
	Key string `envconfig:"SECRET_KEY"`
}
