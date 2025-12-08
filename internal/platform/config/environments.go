// Package config provides centralized configuration structures and loading logic.
package config

import "time"

// GeneralConfig holds a general application configuration.
type GeneralConfig struct {
	Name    string `envconfig:"APP_NAME"`
	Env     string `envconfig:"APP_ENV"`
	Version string `envconfig:"APP_VERSION"`
}

// Secret holds encryption/authentication secrets.
type Secret struct {
	Key string `envconfig:"SECRET_KEY"`
}

// ObservabilityConfig holds all observability-related configuration.
type ObservabilityConfig struct {
	OtelExporterOTLPEndpoint string `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT" default:"otel-collector:4318"`
	OtelServiceName          string `envconfig:"OTEL_SERVICE_NAME"           default:"AionApi"`
	OtelServiceVersion       string `envconfig:"OTEL_SERVICE_VERSION"        default:"0.1.0"`
	OtelExporterHeaders      string `envconfig:"OTEL_EXPORTER_HEADERS"       default:""`
	OtelExporterTimeout      string `envconfig:"OTEL_EXPORTER_TIMEOUT"       default:"5s"`
	OtelExporterCompression  string `envconfig:"OTEL_EXPORTER_COMPRESSION"   default:"none"`
	OtelExporterInsecure     bool   `envconfig:"OTEL_EXPORTER_INSECURE"      default:"true"`
}

// CacheConfig holds Redis cache configuration.
type CacheConfig struct {
	Addr     string `envconfig:"CACHE_ADDR"     default:"redis-aion:6379"`
	Password string `envconfig:"CACHE_PASSWORD"`

	DB             int           `envconfig:"CACHE_REPO"            default:"0"`
	PoolSize       int           `envconfig:"CACHE_POOL_SIZE"       default:"10"`
	ConnectTimeout time.Duration `envconfig:"CACHE_CONNECT_TIMEOUT" default:"5s"`
}

// CookieConfig holds cookie configuration.
type CookieConfig struct {
	Domain   string `envconfig:"COOKIE_DOMAIN"   default:"localhost"`
	Path     string `envconfig:"COOKIE_PATH"     default:"/"`
	SameSite string `envconfig:"COOKIE_SAMESITE" default:"Strict"`
	Secure   bool   `envconfig:"COOKIE_SECURE"   default:"true"`
	MaxAge   int    `envconfig:"COOKIE_MAX_AGE"  default:"0"`
}

// ServerGraphql holds GraphQL server configuration.
type ServerGraphql struct {
	Host string `envconfig:"GRAPHQL_HOST" default:"0.0.0.0"`
	Name string `envconfig:"GRAPHQL_NAME" default:"GraphQL"`
	Path string `envconfig:"GRAPHQL_PATH" default:"/graphql"`

	ReadTimeout  time.Duration `envconfig:"GRAPHQL_READ_TIMEOUT"  default:"5s"`
	WriteTimeout time.Duration `envconfig:"GRAPHQL_WRITE_TIMEOUT" default:"5s"`

	ReadHeaderTimeout time.Duration `envconfig:"GRAPHQL_READ_HEADER_TIMEOUT" default:"5s"`
	IdleTimeout       time.Duration `envconfig:"GRAPHQL_IDLE_TIMEOUT"        default:"60s"`
	MaxHeaderBytes    int           `envconfig:"GRAPHQL_MAX_HEADER_BYTES"    default:"1048576"`
}

// ServerHTTP holds HTTP server configuration.
type ServerHTTP struct {
	Host string `envconfig:"HTTP_HOST" default:"0.0.0.0"`
	Name string `envconfig:"HTTP_NAME" default:"HTTP"`
	Port string `envconfig:"HTTP_PORT" default:"5001"    required:"true"`

	Context          string `envconfig:"HTTP_CONTEXT"            default:"/aion"`
	APIRoot          string `envconfig:"HTTP_API_ROOT"           default:"/api/v1"`
	SwaggerMountPath string `envconfig:"HTTP_SWAGGER_MOUNT_PATH" default:"/swagger"`
	DocsAliasPath    string `envconfig:"HTTP_DOCS_ALIAS_PATH"    default:"/docs"`
	HealthRoute      string `envconfig:"HTTP_HEALTH_ROUTE"       default:"/health"`

	ReadTimeout       time.Duration `envconfig:"HTTP_READ_TIMEOUT"        default:"10s"`
	WriteTimeout      time.Duration `envconfig:"HTTP_WRITE_TIMEOUT"       default:"10s"`
	ReadHeaderTimeout time.Duration `envconfig:"HTTP_READ_HEADER_TIMEOUT" default:"5s"`
	IdleTimeout       time.Duration `envconfig:"HTTP_IDLE_TIMEOUT"        default:"60s"`
	MaxHeaderBytes    int           `envconfig:"HTTP_MAX_HEADER_BYTES"    default:"1048576"` // 1<<20
}

// DBConfig holds database connection configuration.
type DBConfig struct {
	Type            string        `envconfig:"DB_TYPE"                   default:"postgres"`
	Name            string        `envconfig:"DB_NAME"                                       required:"true"`
	Host            string        `envconfig:"DB_HOST"                   default:"localhost"`
	Port            string        `envconfig:"DB_PORT"                   default:"5432"`
	User            string        `envconfig:"DB_USER"                                       required:"true"`
	Password        string        `envconfig:"DB_PASSWORD"                                   required:"true"`
	SSLMode         string        `envconfig:"DB_SSL_MODE"               default:"disable"`
	TimeZone        string        `envconfig:"TIME_ZONE"                 default:"UTC"`
	MaxOpenConns    int           `envconfig:"DB_MAX_CONNECTIONS"        default:"10"`
	MaxIdleConns    int           `envconfig:"DB_MAX_IDLE_CONNECTIONS"   default:"5"`
	MaxRetries      int           `envconfig:"DB_CONNECT_MAX_RETRIES"    default:"3"`
	ConnMaxLifetime time.Duration `envconfig:"DB_CONN_MAX_LIFETIME"      default:"30m"`
	RetryInterval   time.Duration `envconfig:"DB_CONNECT_RETRY_INTERVAL" default:"3s"`
}

// Application holds general application-related configuration.
type Application struct {
	Timeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"5s"`

	ContextRequest time.Duration `envconfig:"CONTEXT_REQUEST" default:"2.1s"`
}

// AionChatConfig holds configuration for the Aion-Chat service (Python AI service).
type AionChatConfig struct {
	BaseURL    string        `envconfig:"AION_CHAT_URL"         default:"http://aion-chat:8000"`
	ServiceKey string        `envconfig:"AION_CHAT_SERVICE_KEY" default:""`
	Timeout    time.Duration `envconfig:"AION_CHAT_TIMEOUT"     default:"30s"`
}
