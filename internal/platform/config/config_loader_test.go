package config_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/stretchr/testify/require"
)

type mockKeyGen struct {
	err error
	key string
}

func (m mockKeyGen) Generate() (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.key, nil
}

type mockConfigLogger struct{}

func (mockConfigLogger) Infof(string, ...any)                      {}
func (mockConfigLogger) Errorf(string, ...any)                     {}
func (mockConfigLogger) Debugf(string, ...any)                     {}
func (mockConfigLogger) Warnf(string, ...any)                      {}
func (mockConfigLogger) Infow(string, ...any)                      {}
func (mockConfigLogger) Errorw(string, ...any)                     {}
func (mockConfigLogger) Debugw(string, ...any)                     {}
func (mockConfigLogger) Warnw(string, ...any)                      {}
func (mockConfigLogger) InfowCtx(context.Context, string, ...any)  {}
func (mockConfigLogger) ErrorwCtx(context.Context, string, ...any) {}
func (mockConfigLogger) WarnwCtx(context.Context, string, ...any)  {}
func (mockConfigLogger) DebugwCtx(context.Context, string, ...any) {}

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DB_NAME", "aion")
	t.Setenv("DB_USER", "aion")
	t.Setenv("DB_PASSWORD", "aion")
	t.Setenv("SETTING_DB_NAME", "aion")
	t.Setenv("SETTING_DB_USER", "aion")
	t.Setenv("SETTING_DB_PASSWORD", "aion")
}

func TestLoaderLoad_GeneratesSecretAndNormalizesContext(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("SETTING_HTTP_CONTEXT", "aion")
	t.Setenv("HTTP_CONTEXT", "aion")
	t.Setenv("SETTING_SECRET_KEY", "")
	t.Setenv("SECRET_KEY", "")

	loader := config.New(mockKeyGen{key: "generated-secret"})
	cfg, err := loader.Load(mockConfigLogger{})
	require.NoError(t, err)
	require.Equal(t, "/aion", cfg.ServerHTTP.Context)
	require.Equal(t, "generated-secret", cfg.Secret.Key)
}

func TestLoaderLoad_UsesExistingSecret(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("SETTING_SECRET_KEY", "already-set")
	t.Setenv("SECRET_KEY", "already-set")

	loader := config.New(mockKeyGen{key: "unused"})
	cfg, err := loader.Load(mockConfigLogger{})
	require.NoError(t, err)
	require.Equal(t, "already-set", cfg.Secret.Key)
}

func TestLoaderLoad_KeyGenerationError(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("SETTING_SECRET_KEY", "")
	t.Setenv("SECRET_KEY", "")

	loader := config.New(mockKeyGen{err: errors.New("keygen failed")})
	_, err := loader.Load(mockConfigLogger{})
	require.Error(t, err)
}
