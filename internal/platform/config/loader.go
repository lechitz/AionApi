package config

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/config/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
)

// Loader is responsible for reading environment configuration.
type Loader struct {
	cfg          Config
	keyGenerator output.KeyGenerator
}

// New returns a new instance of Loader.
func New(keyGen output.KeyGenerator) *Loader {
	return &Loader{
		keyGenerator: keyGen,
	}
}

// Load reads environment configuration and returns a Config struct.
// Returns an error if environment parsing or key generation fails.
func (l *Loader) Load(logger output.ContextLogger) (*Config, error) {
	if err := envconfig.Process(commonkeys.Setting, &l.cfg); err != nil {
		logger.Errorw(constants.ErrFailedToProcessEnvVars, commonkeys.Error, err)
		return nil, err
	}

	if !strings.HasPrefix(l.cfg.ServerHTTP.Context, "/") {
		l.cfg.ServerHTTP.Context = "/" + l.cfg.ServerHTTP.Context
	}

	if l.cfg.Secret.Key == "" {
		generated, err := l.keyGenerator.Generate()
		if err != nil {
			logger.Errorw(constants.ErrGenerateSecretKey, commonkeys.Error, err)
			return nil, err
		}

		l.cfg.Secret.Key = generated

		logger.Warnf(constants.SecretKeyWasNotSet)
		logger.Infof(constants.InfoSecretKeyGenerated, len(generated))
	}

	return &l.cfg, nil
}
