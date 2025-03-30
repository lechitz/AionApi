package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/lechitz/AionApi/internal/infrastructure/security"
	"github.com/lechitz/AionApi/internal/platform/config/constants"
	"github.com/lechitz/AionApi/pkg/utils"
	"go.uber.org/zap"
)

func LoadConfig(logger *zap.SugaredLogger) error {
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
			utils.HandleCriticalError(logger, constants.ErrGenerateSecretKey, err)
			return err
		}
		Setting.SecretKey = generated

		logger.Warn(constants.SecretKeyWasNotSet)
		fmt.Printf("\nSECRET_KEY=%s\n", Setting.SecretKey)
	}

	return nil
}
