package config

import "time"

type DBConfig struct {
	DBUser     string `envconfig:"DB_USER" required:"true"`
	DBPassword string `envconfig:"DB_PASSWORD" required:"true"`
	DBName     string `envconfig:"DB_NAME" required:"true"`
	DBHost     string `envconfig:"DB_HOST" default:"localhost"`
	DBPort     string `envconfig:"DB_PORT" default:"5432"`
	DBType     string `envconfig:"DB_TYPE" default:"postgres"`
}

type CacheConfig struct {
	Addr     string `envconfig:"REDIS_ADDR" default:"redis-aion:6379"`
	Password string `envconfig:"REDIS_PASSWORD"`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
	PoolSize int    `envconfig:"REDIS_POOL_SIZE" default:"10"`
}

var Setting Config

type Config struct {
	Application struct {
		ContextRequest time.Duration `envconfig:"CONTEXT_REQUEST" default:"2.1s"`
	}
	Server struct {
		Context      string        `envconfig:"SERVER_CONTEXT" default:"aion-api"`
		Port         string        `envconfig:"PORT" default:"5001" required:"true"`
		ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"10s"`
		WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
	}
	DBConfig    DBConfig    `envconfig:"POSTGRES"`
	CacheConfig CacheConfig `envconfig:"REDIS"`
	SecretKey   string      `envconfig:"SECRET_KEY"`
}
