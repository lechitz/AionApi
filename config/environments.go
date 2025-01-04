package config

import "time"

var Setting setting

type setting struct {
	Application struct {
		ContextRequest time.Duration `envconfig:"CONTEXT_REQUEST" default:"2.1s"`
	}

	Server struct {
		Context      string        `envconfig:"SERVER_CONTEXT" default:"aion-api"`
		Port         string        `envconfig:"PORT" default:"5001" required:"true" ignored:"false"`
		ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"10s"`
		WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
	}

	Postgres struct {
		DBUser     string `envconfig:"DB_USER" required:"true"`
		DBPassword string `envconfig:"DB_PASSWORD" required:"true"`
		DBName     string `envconfig:"DB_NAME" required:"true"`
		DBHost     string `envconfig:"DB_HOST" default:"localhost"`
		DBPort     string `envconfig:"DB_PORT" default:"5432"`
		DBType     string `envconfig:"DB_TYPE" default:"postgres"`
	}

	SecretKey []byte `envconfig:"SECRET_KEY"`
}
