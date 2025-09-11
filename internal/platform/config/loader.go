package config

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/lechitz/AionApi/internal/platform/ports/output/keygen"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
)

// Loader is responsible for reading environment configuration.
type Loader struct {
	cfg          Config
	keyGenerator keygen.Generator
}

// New returns a new instance of Loader.
func New(keyGen keygen.Generator) *Loader {
	return &Loader{
		keyGenerator: keyGen,
	}
}

// Load reads environment configuration and returns a Config struct.
// Returns an error if environment parsing or key generation fails.
func (l *Loader) Load(logger logger.ContextLogger) (*Config, error) {
	if err := envconfig.Process(commonkeys.Setting, &l.cfg); err != nil {
		logger.Errorw(ErrFailedToProcessEnvVars, commonkeys.Error, err)
		return nil, err
	}

	if !strings.HasPrefix(l.cfg.ServerHTTP.Context, "/") {
		l.cfg.ServerHTTP.Context = "/" + l.cfg.ServerHTTP.Context
	}

	if l.cfg.Secret.Key == "" {
		generated, err := l.keyGenerator.Generate()
		if err != nil {
			logger.Errorw(ErrGenerateSecretKey, commonkeys.Error, err)
			return nil, err
		}

		l.cfg.Secret.Key = generated

		logger.Warnf(SecretKeyWasNotSet)
		logger.Infof(InfoSecretKeyGenerated, len(generated))
	}

	return &l.cfg, nil
}
