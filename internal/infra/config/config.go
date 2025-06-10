package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/adapters/secondary/security"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/infra/config/constants"
)

// Load initializes application settings by processing environment variables and generating a default secret key if not set. Logs warnings or errors as necessary.
func Load(logger logger.Logger) error {
	if err := envconfig.Process(constants.Settings, &Setting); err != nil {
		response.HandleCriticalError(logger, constants.ErrFailedToProcessEnvVars, err)
	}

	if Setting.Secret.Key == "" {
		generated, err := security.GenerateJWTKey()
		if err != nil {
			response.HandleCriticalError(logger, constants.ErrGenerateSecretKey, err)
		}

		Setting.Secret.Key = generated

		logger.Warnf(constants.SecretKeyWasNotSet)
		logger.Infof("JWT secret key successfully generated with length: %d", len(generated))
	}

	return nil
}
