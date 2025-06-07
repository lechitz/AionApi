package config

import "time"

var Setting Config

type Config struct {
	Application   Application
	ServerHTTP    ServerHTTP
	ServerGraphql ServerGraphql
	DB            DBConfig
	Cache         CacheConfig
	Secret        Secret
}

type DBConfig struct {
	Name     string `envconfig:"DB_NAME"     required:"true"`
	Host     string `envconfig:"DB_HOST"                     default:"localhost"`
	Port     string `envconfig:"DB_PORT"                     default:"5432"`
	Type     string `envconfig:"DB_TYPE"                     default:"postgres"`
	User     string `envconfig:"DB_USER"     required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
}

type CacheConfig struct {
	Addr     string `envconfig:"CACHE_ADDR"      default:"redis-aion:6379"`
	Password string `envconfig:"CACHE_PASSWORD"`
	DB       int    `envconfig:"CACHE_REPO"      default:"0"`
	PoolSize int    `envconfig:"CACHE_POOL_SIZE" default:"10"`
}

type ServerGraphql struct {
	Port string `envconfig:"GRAPHQL_PORT" default:"8081" required:"true"`
}

type ServerHTTP struct {
	Context      string        `envconfig:"HTTP_CONTEXT"       default:"aion-api"`
	Port         string        `envconfig:"HTTP_PORT"          default:"5001"     required:"true"`
	ReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT"  default:"10s"`
	WriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"10s"`
}

type Application struct {
	ContextRequest time.Duration `envconfig:"CONTEXT_REQUEST"  default:"2.1s"`
	Timeout        int           `envconfig:"SHUTDOWN_TIMEOUT" default:"5"`
}

type Secret struct {
	Key string `envconfig:"SECRET_KEY"`
}
