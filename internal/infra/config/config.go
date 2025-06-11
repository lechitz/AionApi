// Package config provides configuration management for the application.
package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/internal/adapters/secondary/security"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/infra/config/constants"
)

// Loader is responsible for reading environment configuration
// and returning a fully populated Config object.
// This struct helps avoid global state and improves testability.
type Loader struct {
	cfg Config
}

// NewLoader returns a new instance of Loader.
// This constructor enforces explicit configuration loading.
func NewLoader() *Loader {
	return &Loader{}
}

// Load reads configuration from environment variables into Config.
// It also generates a default JWT secret key if not explicitly set.
// Parameters:
//   - logger: application-wide logger used for logging critical failures.
//
// Returns:
//   - Config: fully loaded configuration object
//   - error: error occurred during loading or key generation
func (l *Loader) Load(logger logger.Logger) (Config, error) {
	if err := envconfig.Process(constants.Settings, &l.cfg); err != nil {
		response.HandleCriticalError(logger, constants.ErrFailedToProcessEnvVars, err)
		return Config{}, err
	}

	if l.cfg.Secret.Key == "" {
		generated, err := security.GenerateJWTKey()
		if err != nil {
			response.HandleCriticalError(logger, constants.ErrGenerateSecretKey, err)
			return Config{}, err
		}

		l.cfg.Secret.Key = generated
		logger.Warnf(constants.SecretKeyWasNotSet)
		logger.Infof("JWT secret key successfully generated with length: %d", len(generated))
	}

	return l.cfg, nil
}

// Config holds all configuration sections required to bootstrap the application.
// This struct is populated through environment variable processing via envconfig.
type Config struct {
	DB            DBConfig      // Database-related configuration
	Cache         CacheConfig   // Cache layer configuration
	Secret        Secret        // Secret and security configuration
	ServerGraphql ServerGraphql // GraphQL server-specific configuration
	ServerHTTP    ServerHTTP    // HTTP server-specific configuration
	Application   Application   // General application-level configuration
}
