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

// loggerSuggar pode ser injetado ou setado em outro lugar, caso queira usar.
// Aqui, deixamos apenas como referência se precisar de logs de erro críticos.
var loggerSuggar *zap.SugaredLogger

func LoadConfig() error {
	// 1) Processa as variáveis de ambiente e joga em Setting
	err := envconfig.Process(contextkeys.Setting, &Setting)
	if err != nil {
		return fmt.Errorf(msg.ErrFailedToProcessEnvVars, err)
	}

	// 2) Validações simples
	if Setting.Server.Context == "" {
		return fmt.Errorf(msg.ErrServerContextEmpty)
	}
	if Setting.Server.Port == "" {
		return fmt.Errorf(msg.ErrServerPortEmpty)
	}

	// 3) Se SECRET_KEY estiver vazia, gera uma nova
	if len(Setting.SecretKey) == 0 {
		generated, err := utils.GenerateJWTKey()
		if err != nil {
			errors.HandleCriticalError(loggerSuggar, msg.ErrGenerateSecretKey, err)
			return err
		}
		// Atribui a Setting.SecretKey
		Setting.SecretKey = generated

		// Imprime no console para o usuário ver e copiar
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
