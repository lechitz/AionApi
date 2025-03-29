package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/lechitz/AionApi/core/msg"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"github.com/lechitz/AionApi/pkg/errors"
	"github.com/lechitz/AionApi/pkg/utils"
	"go.uber.org/zap"
)

func LoadConfig(logger *zap.SugaredLogger) error {
	if err := envconfig.Process(contextkeys.Setting, &Setting); err != nil {
		return fmt.Errorf(msg.ErrFailedToProcessEnvVars, err)
	}

	if Setting.Server.Context == "" {
		return fmt.Errorf(msg.ErrServerContextEmpty)
	}
	if Setting.Server.Port == "" {
		return fmt.Errorf(msg.ErrServerPortEmpty)
	}

	if Setting.SecretKey == "" {
		generated, err := utils.GenerateJWTKey()
		if err != nil {
			errors.HandleCriticalError(logger, msg.ErrGenerateSecretKey, err)
			return err
		}
		Setting.SecretKey = generated

		logger.Warn("SECRET_KEY was not set. A new one was generated for this runtime session.")
		fmt.Printf("\n👉 SECRET_KEY=%s\n", Setting.SecretKey)
	}

	return nil
}
