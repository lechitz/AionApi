package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/adapters/secondary/security"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/infra/config/constants"
)

// setting holds the application-wide configuration after loading from environment.
var setting Config

// Get returns the loaded Config.
func Get() Config {
	return setting
}

// Setting returns a copy of the loaded application configuration.
func Setting() *Config {
	return &setting
}

// Config holds all configuration sections for initializing the application.
type Config struct {
	DB            DBConfig
	Cache         CacheConfig
	Secret        Secret
	ServerGraphql ServerGraphql
	ServerHTTP    ServerHTTP
	Application   Application
}

// Load loads configuration from environment into Setting.
func Load(logger logger.Logger) error {
	if err := envconfig.Process(constants.Settings, &setting); err != nil {
		response.HandleCriticalError(logger, constants.ErrFailedToProcessEnvVars, err)
	}

	if setting.Secret.Key == "" {
		generated, err := security.GenerateJWTKey()
		if err != nil {
			response.HandleCriticalError(logger, constants.ErrGenerateSecretKey, err)
		}

		setting.Secret.Key = generated
		logger.Warnf(constants.SecretKeyWasNotSet)
		logger.Infof("JWT secret key successfully generated with length: %d", len(generated))
	}

	return nil
}
