// Package config provides configuration management for the application.
package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/lechitz/AionApi/internal/shared/common"

	"github.com/kelseyhightower/envconfig"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/internal/adapters/secondary/security"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/config/constants"
)

// Loader is responsible for reading environment configuration.
type Loader struct {
	cfg Config
}

// NewLoader returns a new instance of Loader.
func NewLoader() *Loader {
	return &Loader{}
}

// Load reads configuration from environment variables into Config.
func (l *Loader) Load(logger output.Logger) (*Config, error) {
	if err := envconfig.Process(common.Setting, &l.cfg); err != nil {
		response.HandleCriticalError(logger, constants.ErrFailedToProcessEnvVars, err) // TODO: AJUSTAR ERRO
		return nil, err
	}

	if l.cfg.Secret.Key == "" {
		generated, err := security.GenerateJWTKey()
		if err != nil {
			response.HandleCriticalError(logger, constants.ErrGenerateSecretKey, err) // TODO: AJUSTAR ERRO
			return nil, err
		}

		l.cfg.Secret.Key = generated
		logger.Warnf(constants.SecretKeyWasNotSet)
		logger.Infof("JWT secret key successfully generated with length: %d", len(generated)) // TODO: AJUSTAR LOGGER
	}

	return &l.cfg, nil
}

// Config holds all configuration sections required to bootstrap the application.
// This struct is populated through environment variable processing via envconfig.
type Config struct {
	Secret        Secret
	Observability ObservabilityConfig
	ServerGraphql ServerGraphql
	ServerHTTP    ServerHTTP
	Cache         CacheConfig
	DB            DBConfig
	Application   Application
}

// Validate checks if all mandatory fields are set and correct.
func (c *Config) Validate() error {
	if err := c.validateHTTP(); err != nil {
		return err
	}
	if err := c.validateGraphQL(); err != nil {
		return err
	}
	if err := c.validateCache(); err != nil {
		return err
	}
	if err := c.validateDB(); err != nil {
		return err
	}
	if err := c.validateObservability(); err != nil {
		return err
	}
	if err := c.validateApp(); err != nil {
		return err
	}

	return nil
}

// TODO: AJUSTAR ERROS ABAIXO EM CONST !

func (c *Config) validateHTTP() error {
	if c.ServerHTTP.Port == "" {
		return errors.New("HTTP port is required")
	}
	if c.ServerHTTP.ReadTimeout < time.Second { // TODO: Passar o valor para variavel de ambiente
		return errors.New("HTTP read timeout must be at least 1s")
	} // TODO: avaliar se não deveria ir para variáveis de ambiente
	if c.ServerHTTP.WriteTimeout < time.Second { // TODO: Passar o valor para variavel de ambiente
		return errors.New("HTTP write timeout must be at least 1s")
	} // TODO: avaliar se não deveria ir para variáveis de ambiente
	return nil
}

func (c *Config) validateGraphQL() error {
	if c.ServerGraphql.Port == "" {
		return errors.New("GraphQL port is required")
	}
	if c.ServerGraphql.Path == "" {
		return errors.New("GraphQL path is required")
	}
	if c.ServerGraphql.Path[0] != '/' { // TODO: Passar o valor para variavel de ambiente
		return errors.New("GraphQL path must start with '/'")
	} // TODO: avaliar se não deveria ir para variáveis de ambiente
	if c.ServerGraphql.ReadTimeout < time.Second { // TODO: Passar o valor para variavel de ambiente
		return errors.New("GraphQL read timeout must be at least 1s")
	} // TODO: avaliar se não deveria ir para variáveis de ambiente
	if c.ServerGraphql.WriteTimeout < time.Second { // TODO: Passar o valor para variavel de ambiente
		return errors.New("GraphQL write timeout must be at least 1s")
	} // TODO: avaliar se não deveria ir para variáveis de ambiente
	return nil
}

func (c *Config) validateCache() error {
	if c.Cache.PoolSize < 1 {
		return errors.New("cache pool size must be at least 1")
	} // TODO: avaliar se não deveria ir para variáveis de ambiente
	if c.Cache.Addr == "" {
		return errors.New("cache address cannot be empty")
	}
	return nil
}

func (c *Config) validateDB() error {
	if c.DB.Name == "" {
		return errors.New("database name is required")
	}
	if c.DB.User == "" {
		return errors.New("database user is required")
	}
	if c.DB.Password == "" {
		return errors.New("database password is required")
	}
	if c.DB.MaxOpenConns < 1 { // TODO: Passar o valor para variavel de ambiente
		return errors.New("DB_MAX_CONNECTIONS must be at least 1")
	}
	if c.DB.MaxIdleConns < 0 { // TODO: Passar o valor para variavel de ambiente
		return errors.New("DB_MAX_IDLE_CONNECTIONS cannot be negative")
	}
	if c.DB.ConnMaxLifetime < 0 { // TODO: Passar o valor para variavel de ambiente
		return errors.New("DB_CONN_MAX_LIFETIME_MINUTES cannot be negative")
	}
	return nil
}

func (c *Config) validateObservability() error {
	if c.Observability.OtelExporterOTLPEndpoint == "" {
		return errors.New("OTel Exporter endpoint cannot be empty")
	}
	if c.Observability.OtelExporterCompression != "" {
		switch c.Observability.OtelExporterCompression {
		case "none", "gzip": // TODO: avaliar o uso.
			// ok
		default:
			return fmt.Errorf("OTel Exporter compression must be either 'none' or 'gzip', got: %s", c.Observability.OtelExporterCompression)
		}
	}
	return nil
}

func (c *Config) validateApp() error {
	if c.Application.ContextRequest < 500*time.Millisecond { // TODO: Passar o valor para variavel de ambiente
		return errors.New("context request timeout must be at least 500ms")
	} // TODO: avaliar se não deveria ir para variáveis de ambiente
	if c.Application.Timeout < 1 { // TODO: Passar o valor para variavel de ambiente
		return errors.New("shutdown timeout must be at least 1 second")
	} // TODO: avaliar se não deveria ir para variáveis de ambiente
	return nil
}
