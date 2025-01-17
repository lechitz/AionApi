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

var loggerSuggar *zap.SugaredLogger

func LoadConfig() error {

	err := envconfig.Process(contextkeys.Setting, &Setting)
	if err != nil {
		return fmt.Errorf(msg.ErrFailedToProcessEnvVars, err)
	}

	if Setting.Server.Context == "" {
		return fmt.Errorf(msg.ErrServerContextEmpty)
	}
	if Setting.Server.Port == "" {
		return fmt.Errorf(msg.ErrServerPortEmpty)
	}

	if len(Setting.SecretKey) == 0 {
		generated, err := utils.GenerateJWTKey()
		if err != nil {
			errors.HandleCriticalError(loggerSuggar, msg.ErrGenerateSecretKey, err)
			return err
		}
		Setting.SecretKey = generated

		fmt.Println()
		fmt.Println("====================================================================")
		fmt.Println("SECRET_KEY was not found. A new one has been generated.")
		fmt.Println("Please copy and add this to your .env if you want to reuse it later:")
		fmt.Println(" ")
		fmt.Printf("SECRET_KEY=%s\n", Setting.SecretKey)
		fmt.Println("====================================================================")
		fmt.Println()
	}

	return nil
}
