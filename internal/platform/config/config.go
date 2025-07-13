// Package config provides configuration loading and validation for the application.
package config

import (
	"errors"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/lechitz/AionApi/internal/adapters/secondary/security"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/config/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
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
// Panics on fatal error, logging with the provided logger.
func (l *Loader) Load(logger output.ContextLogger) (*Config, error) {
	if err := envconfig.Process(commonkeys.Setting, &l.cfg); err != nil {
		logger.Errorw(constants.ErrFailedToProcessEnvVars, commonkeys.Error, err)
		panic(err)
	}

	if l.cfg.Secret.Key == "" {
		generated, err := security.GenerateJWTKey()
		if err != nil {
			logger.Errorw(constants.ErrGenerateSecretKey, commonkeys.Error, err)
			panic(err)
		}
		l.cfg.Secret.Key = generated
		logger.Warnf(constants.SecretKeyWasNotSet)
		logger.Infof(constants.InfoJWTSecretGenerated, len(generated))
	}

	return &l.cfg, nil
}

// Validate checks if the configuration is valid, returning the first validation error encountered.
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
		return errors.New(constants.ErrHTTPPortRequired)
	}
	if c.ServerHTTP.Context == "" {
		return errors.New(constants.ErrHTTPContextPathEmpty)
	}
	if c.ServerHTTP.ReadTimeout < constants.MinHTTPTimeout {
		return fmt.Errorf(constants.ErrHTTPReadTimeoutMin, constants.MinHTTPTimeout)
	}
	if c.ServerHTTP.WriteTimeout < constants.MinHTTPTimeout {
		return fmt.Errorf(constants.ErrHTTPWriteTimeoutMin, constants.MinHTTPTimeout)
	}
	return nil
}

func (c *Config) validateGraphQL() error {
	if c.ServerGraphql.Port == "" {
		return errors.New(constants.ErrGraphqlPortRequired)
	}
	if c.ServerGraphql.Path == "" {
		return errors.New(constants.ErrGraphqlPathRequired)
	}
	if c.ServerGraphql.Path[0] != '/' {
		return errors.New(constants.ErrGraphqlPathMustStart)
	}
	if c.ServerGraphql.ReadTimeout < constants.MinGraphQLTimeout {
		return fmt.Errorf(constants.ErrGraphqlReadTimeoutMin, constants.MinGraphQLTimeout)
	}
	if c.ServerGraphql.WriteTimeout < constants.MinGraphQLTimeout {
		return fmt.Errorf(constants.ErrGraphqlWriteTimeoutMin, constants.MinGraphQLTimeout)
	}
	return nil
}

func (c *Config) validateCache() error {
	if c.Cache.PoolSize < constants.MinCachePoolSize {
		return fmt.Errorf(constants.ErrCachePoolSizeMin, constants.MinCachePoolSize)
	}
	if c.Cache.Addr == "" {
		return errors.New(constants.ErrCacheAddrEmpty)
	}
	return nil
}

func (c *Config) validateDB() error {
	if c.DB.Type == "" {
		return errors.New(constants.ErrDBTypeEmpty)
	}
	if c.DB.Type != "postgres" {
		return fmt.Errorf(constants.ErrDBTypeUnsupported, c.DB.Type)
	}
	if c.DB.Host == "" {
		return errors.New(constants.ErrDBHostEmpty)
	}
	if c.DB.Port == "" {
		return errors.New(constants.ErrDBPortEmpty)
	}
	if c.DB.Name == "" {
		return errors.New(constants.ErrDBNameRequired)
	}
	if c.DB.User == "" {
		return errors.New(constants.ErrDBUserRequired)
	}
	if c.DB.Password == "" {
		return errors.New(constants.ErrDBPasswordRequired)
	}
	if c.DB.TimeZone == "" {
		return errors.New(constants.ErrDBTimeZoneEmpty)
	}

	switch c.DB.SSLMode {
	case "disable", "require", "verify-ca", "verify-full":
		// valid
	default:
		return fmt.Errorf(constants.ErrDBSSLModeInvalid, c.DB.SSLMode)
	}

	if c.DB.MaxOpenConns < constants.MinDBMaxOpenConns {
		return fmt.Errorf(constants.ErrDBMaxOpenConnsMin, constants.MinDBMaxOpenConns)
	}
	if c.DB.MaxIdleConns < constants.MinDBMaxIdleConns {
		return errors.New(constants.ErrDBMaxIdleConnsNeg)
	}
	if c.DB.ConnMaxLifetime < constants.MinDBConnMaxLifetimeMin {
		return errors.New(constants.ErrDBConnMaxLifetimeNeg)
	}
	if c.DB.RetryInterval < constants.MinDBRetryInterval {
		return fmt.Errorf(constants.ErrDBRetryIntervalMin, constants.MinDBRetryInterval)
	}
	if c.DB.MaxRetries < constants.MinDBMaxRetries {
		return fmt.Errorf(constants.ErrDBMaxRetriesMin, constants.MinDBMaxRetries)
	}
	return nil
}

func (c *Config) validateObservability() error {
	if c.Observability.OtelExporterOTLPEndpoint == "" {
		return errors.New(constants.ErrOtelEndpointEmpty)
	}
	if c.Observability.OtelExporterCompression != "" {
		switch c.Observability.OtelExporterCompression {
		case "none", "gzip":
			// valid
		default:
			return fmt.Errorf(
				constants.ErrOtelCompressionInvalid,
				c.Observability.OtelExporterCompression,
			)
		}
	}
	return nil
}

func (c *Config) validateApp() error {
	if c.Application.ContextRequest < constants.MinContextRequest {
		return fmt.Errorf(constants.ErrAppContextReqMin, constants.MinContextRequest)
	}
	if c.Application.Timeout < constants.MinShutdownTimeout {
		return fmt.Errorf(constants.ErrAppShutdownTimeoutMin, constants.MinShutdownTimeout)
	}
	return nil
}
