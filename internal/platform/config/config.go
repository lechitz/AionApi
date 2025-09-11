// Package config provides configuration loading and validation for the application.
package config

import (
	"errors"
	"fmt"
)

// Config holds all configuration sections required to bootstrap the application.
type Config struct {
	General       GeneralConfig
	Secret        Secret
	Observability ObservabilityConfig
	Cache         CacheConfig
	Cookie        CookieConfig
	ServerGraphql ServerGraphql
	ServerHTTP    ServerHTTP
	ServerGRPC    ServerGRPC
	DB            DBConfig
	Application   Application
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
	if c.ServerHTTP.Host == "" {
		return errors.New(ErrHTTPHostRequired)
	}
	if c.ServerHTTP.Port == "" {
		return errors.New(ErrHTTPPortRequired)
	}
	if c.ServerHTTP.Context == "" {
		return errors.New(ErrHTTPContextPathEmpty)
	}
	if c.ServerHTTP.Context[0] != '/' {
		return errors.New(ErrHTTPContextMustStart)
	}
	if len(c.ServerHTTP.Context) > 1 && c.ServerHTTP.Context[len(c.ServerHTTP.Context)-1] == '/' {
		return errors.New(ErrHTTPContextMustNotEndWithSlash)
	}
	if c.ServerHTTP.ReadTimeout < MinHTTPTimeout {
		return fmt.Errorf(ErrHTTPReadTimeoutMin, MinHTTPTimeout)
	}
	if c.ServerHTTP.WriteTimeout < MinHTTPTimeout {
		return fmt.Errorf(ErrHTTPWriteTimeoutMin, MinHTTPTimeout)
	}
	if c.ServerHTTP.ReadHeaderTimeout <= 0 {
		return errors.New(ErrHTTPReadHeaderTimeoutMin)
	}
	if c.ServerHTTP.IdleTimeout <= 0 {
		return errors.New(ErrHTTPIdleTimeoutMin)
	}
	if c.ServerHTTP.MaxHeaderBytes <= 0 {
		return errors.New(ErrHTTPMaxHeaderBytesMin)
	}

	return nil
}

func (c *Config) validateGraphQL() error {
	if c.ServerGraphql.Path == "" {
		return errors.New(ErrGraphqlPathRequired)
	}
	if c.ServerGraphql.Path[0] != '/' {
		return errors.New(ErrGraphqlPathMustStart)
	}
	//if c.ServerGraphql.ReadTimeout < MinGraphQLTimeout {
	//	return fmt.Errorf(ErrGraphqlReadTimeoutMin, MinGraphQLTimeout)
	//}
	//if c.ServerGraphql.WriteTimeout < MinGraphQLTimeout {
	//	return fmt.Errorf(ErrGraphqlWriteTimeoutMin, MinGraphQLTimeout)
	//}
	//if c.ServerGraphql.ReadHeaderTimeout <= 0 {
	//	return errors.New(ErrGraphqlReadHeaderTimeoutMin)
	//}
	//if c.ServerGraphql.IdleTimeout <= 0 {
	//	return errors.New(ErrHTTPIdleTimeoutMin)
	//}
	//if c.ServerGraphql.MaxHeaderBytes <= 0 {
	//	return errors.New(ErrHTTPMaxHeaderBytesMin)
	//}

	return nil
}

func (c *Config) validateCache() error {
	if c.Cache.PoolSize < MinCachePoolSize {
		return fmt.Errorf(ErrCachePoolSizeMin, MinCachePoolSize)
	}
	if c.Cache.Addr == "" {
		return errors.New(ErrCacheAddrEmpty)
	}
	return nil
}

func (c *Config) validateDB() error {
	if c.DB.Type == "" {
		return errors.New(ErrDBTypeEmpty)
	}
	if c.DB.Host == "" {
		return errors.New(ErrDBHostEmpty)
	}
	if c.DB.Port == "" {
		return errors.New(ErrDBPortEmpty)
	}
	if c.DB.Name == "" {
		return errors.New(ErrDBNameRequired)
	}
	if c.DB.User == "" {
		return errors.New(ErrDBUserRequired)
	}
	if c.DB.Password == "" {
		return errors.New(ErrDBPasswordRequired)
	}
	if c.DB.TimeZone == "" {
		return errors.New(ErrDBTimeZoneEmpty)
	}

	switch c.DB.SSLMode {
	case "disable", "require", "verify-ca", "verify-full":
		// valid
	default:
		return fmt.Errorf(ErrDBSSLModeInvalid, c.DB.SSLMode)
	}

	if c.DB.MaxOpenConns < MinDBMaxOpenConns {
		return fmt.Errorf(ErrDBMaxOpenConnsMin, MinDBMaxOpenConns)
	}
	if c.DB.MaxIdleConns < MinDBMaxIdleConns {
		return errors.New(ErrDBMaxIdleConnsNeg)
	}
	if c.DB.ConnMaxLifetime < MinDBConnMaxLifetimeMin {
		return errors.New(ErrDBConnMaxLifetimeNeg)
	}
	if c.DB.RetryInterval < MinDBRetryInterval {
		return fmt.Errorf(ErrDBRetryIntervalMin, MinDBRetryInterval)
	}
	if c.DB.MaxRetries < MinDBMaxRetries {
		return fmt.Errorf(ErrDBMaxRetriesMin, MinDBMaxRetries)
	}

	return nil
}

func (c *Config) validateObservability() error {
	if c.Observability.OtelExporterOTLPEndpoint == "" {
		return errors.New(ErrOtelEndpointEmpty)
	}
	if c.Observability.OtelExporterCompression != "" {
		switch c.Observability.OtelExporterCompression {
		case "none", "gzip":
			// valid
		default:
			return fmt.Errorf(
				ErrOtelCompressionInvalid,
				c.Observability.OtelExporterCompression,
			)
		}
	}

	return nil
}

func (c *Config) validateApp() error {
	if c.Application.ContextRequest < MinContextRequest {
		return fmt.Errorf(ErrAppContextReqMin, MinContextRequest)
	}
	if c.Application.Timeout < MinShutdownTimeout {
		return fmt.Errorf(ErrAppShutdownTimeoutMin, MinShutdownTimeout)
	}

	return nil
}
