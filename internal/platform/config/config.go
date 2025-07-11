package config

import (
	"errors"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/internal/adapters/secondary/security"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/config/constants"
	"github.com/lechitz/AionApi/internal/shared/commonkeys"
)

// Loader is responsible for reading environment configuration.
type Loader struct {
	cfg Config
}

// New returns a new instance of Loader.
func New() *Loader {
	return &Loader{}
}

// Config holds all configuration sections required to bootstrap the application.
type Config struct {
	General       GeneralConfig
	Secret        Secret
	Observability ObservabilityConfig
	Cache         CacheConfig
	Cookie        CookieConfig
	ServerGraphql ServerGraphql
	ServerHTTP    ServerHTTP
	DB            DBConfig
	Application   Application
}

// Load reads environment configuration and returns a Config struct.
func (l *Loader) Load(logger output.ContextLogger) (*Config, error) {
	if err := envconfig.Process(commonkeys.Setting, &l.cfg); err != nil {
		response.HandleCriticalError(logger, constants.ErrFailedToProcessEnvVars, err)
		return nil, err
	}

	if l.cfg.Secret.Key == "" {
		generated, err := security.GenerateJWTKey()
		if err != nil {
			response.HandleCriticalError(logger, constants.ErrGenerateSecretKey, err)
			return nil, err
		}
		l.cfg.Secret.Key = generated
		logger.Warnf(constants.SecretKeyWasNotSet)
		logger.Infof("JWT secret key successfully generated with length: %d", len(generated)) // TODO: ajustar magic string.
	}

	return &l.cfg, nil
}

// Validate checks if the configuration is valid.
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

func (c *Config) validateHTTP() error {
	if c.ServerHTTP.Port == "" {
		return errors.New("HTTP port is required") // TODO: ajustar magic string.
	}
	if c.ServerHTTP.Context == "" {
		return errors.New("HTTP context path cannot be empty") // TODO: ajustar magic string.
	}
	if c.ServerHTTP.ReadTimeout < constants.MinHTTPTimeout {
		return fmt.Errorf("HTTP read timeout must be at least %v", constants.MinHTTPTimeout) // TODO: ajustar magic string.
	}
	if c.ServerHTTP.WriteTimeout < constants.MinHTTPTimeout {
		return fmt.Errorf("HTTP write timeout must be at least %v", constants.MinHTTPTimeout) // TODO: ajustar magic string.
	}
	return nil
}

func (c *Config) validateGraphQL() error {
	if c.ServerGraphql.Port == "" {
		return errors.New("GraphQL port is required") // TODO: ajustar magic string.
	}
	if c.ServerGraphql.Path == "" {
		return errors.New("GraphQL path is required") // TODO: ajustar magic string.
	}
	if c.ServerGraphql.Path[0] != '/' { // TODO: ajustar magic string.
		return errors.New("GraphQL path must start with '/'") // TODO: ajustar magic string.
	}
	if c.ServerGraphql.ReadTimeout < constants.MinGraphQLTimeout {
		return fmt.Errorf("GraphQL read timeout must be at least %v", constants.MinGraphQLTimeout) // TODO: ajustar magic string.
	}
	if c.ServerGraphql.WriteTimeout < constants.MinGraphQLTimeout {
		return fmt.Errorf("GraphQL write timeout must be at least %v", constants.MinGraphQLTimeout) // TODO: ajustar magic string.
	}
	return nil
}

func (c *Config) validateCache() error {
	if c.Cache.PoolSize < constants.MinCachePoolSize {
		return fmt.Errorf("CACHE_POOL_SIZE must be at least %d", constants.MinCachePoolSize) // TODO: ajustar magic string.
	}
	if c.Cache.Addr == "" {
		return errors.New("cache address cannot be empty") // TODO: ajustar magic string.
	}
	return nil
}

func (c *Config) validateDB() error {
	if c.DB.Type == "" {
		return errors.New("DB_TYPE cannot be empty") // TODO: ajustar magic string.
	}
	if c.DB.Type != "postgres" { // TODO: ajustar magic string.
		return fmt.Errorf("unsupported DB_TYPE: %s", c.DB.Type) // TODO: ajustar magic string.
	}
	if c.DB.Host == "" {
		return errors.New("DB_HOST cannot be empty") // TODO: ajustar magic string.
	}
	if c.DB.Port == "" {
		return errors.New("DB_PORT cannot be empty") // TODO: ajustar magic string.
	}
	if c.DB.Name == "" {
		return errors.New("database name is required") // TODO: ajustar magic string.
	}
	if c.DB.User == "" {
		return errors.New("database user is required") // TODO: ajustar magic string.
	}
	if c.DB.Password == "" {
		return errors.New("database password is required") // TODO: ajustar magic string.
	}
	if c.DB.TimeZone == "" {
		return errors.New("DB timezone (TZ) cannot be empty") // TODO: ajustar magic string.
	}

	switch c.DB.SSLMode {
	case "disable", "require", "verify-ca", "verify-full": // TODO: ajustar magic string.
		// Next t
	default:
		return fmt.Errorf("invalid DB_SSL_MODE value: %s", c.DB.SSLMode) // TODO: ajustar magic string.
	}

	if c.DB.MaxOpenConns < constants.MinDBMaxOpenConns {
		return fmt.Errorf("DB_MAX_CONNECTIONS must be at least %d", constants.MinDBMaxOpenConns) // TODO: ajustar magic string.
	}
	if c.DB.MaxIdleConns < constants.MinDBMaxIdleConns {
		return errors.New("DB_MAX_IDLE_CONNECTIONS cannot be negative") // TODO: ajustar magic string.
	}
	if c.DB.ConnMaxLifetime < constants.MinDBConnMaxLifetimeMin {
		return errors.New("DB_CONN_MAX_LIFETIME_MINUTES cannot be negative") // TODO: ajustar magic string.
	}
	if c.DB.RetryInterval < constants.MinDBRetryInterval {
		return fmt.Errorf("DB_CONNECT_RETRY_INTERVAL must be at least %v", constants.MinDBRetryInterval) // TODO: ajustar magic string.
	}
	if c.DB.MaxRetries < constants.MinDBMaxRetries {
		return fmt.Errorf("DB_CONNECT_MAX_RETRIES must be at least %d", constants.MinDBMaxRetries) // TODO: ajustar magic string.
	}
	return nil
}

func (c *Config) validateObservability() error {
	if c.Observability.OtelExporterOTLPEndpoint == "" {
		return errors.New("OTel Exporter endpoint cannot be empty") // TODO: ajustar magic string.
	}
	if c.Observability.OtelExporterCompression != "" {
		switch c.Observability.OtelExporterCompression {
		case "none", "gzip": // TODO: ajustar magic string.
			// ok
		default:
			return fmt.Errorf(
				"OTel Exporter compression must be either 'none' or 'gzip', got: %s",
				c.Observability.OtelExporterCompression,
			) // TODO: ajustar magic string.
		}
	}
	return nil
}

func (c *Config) validateApp() error {
	if c.Application.ContextRequest < constants.MinContextRequest {
		return fmt.Errorf("context request timeout must be at least %v", constants.MinContextRequest) // TODO: ajustar magic string.
	}
	if c.Application.Timeout < constants.MinShutdownTimeout {
		return fmt.Errorf("shutdown timeout must be at least %d second", constants.MinShutdownTimeout) // TODO: ajustar magic string.
	}
	return nil
}
