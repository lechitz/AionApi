package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/adapters/secondary/security"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/infra/config/constants"
)

func Load(logger logger.Logger) error {
	if err := envconfig.Process(constants.Settings, &Setting); err != nil {
		return fmt.Errorf(constants.ErrFailedToProcessEnvVars, err)
	}

	if Setting.Server.Context == "" {
		return fmt.Errorf(constants.ErrServerContextEmpty)
	}
	if Setting.Server.Port == "" {
		return fmt.Errorf(constants.ErrServerPortEmpty)
	}

	if Setting.SecretKey == "" {
		generated, err := security.GenerateJWTKey()
		if err != nil {
			response.HandleCriticalError(logger, constants.ErrGenerateSecretKey, err)
			return err
		}

		Setting.SecretKey = generated

		logger.Warnf(constants.SecretKeyWasNotSet)
		fmt.Printf(constants.SecretKeyFormat, Setting.SecretKey)
	}

	return nil
}
