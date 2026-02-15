//nolint:testpackage // tests exercise package-private wiring helpers.
package fxapp

import (
	"errors"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/tests/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setRequiredSettingEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DB_NAME", "aion")
	t.Setenv("DB_USER", "aion")
	t.Setenv("DB_PASSWORD", "aion")
	t.Setenv("SETTING_DB_NAME", "aion")
	t.Setenv("SETTING_DB_USER", "aion")
	t.Setenv("SETTING_DB_PASSWORD", "aion")
	t.Setenv("SETTING_SECRET_KEY", "already-set")
}

func TestProvideLoggerRegistersOnStopHook(t *testing.T) {
	lc := &fakeLifecycle{}

	log := ProvideLogger(lc)
	require.NotNil(t, log)
	require.Len(t, lc.hooks, 1)
	require.NotNil(t, lc.hooks[0].OnStop)
	require.NoError(t, lc.hooks[0].OnStop(t.Context()))
}

func TestProvideConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		setRequiredSettingEnv(t)
		t.Setenv("SECRET_KEY", "")
		t.Setenv("SETTING_SECRET_KEY", "")

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		keyGen := mocks.NewMockGenerator(ctrl)
		keyGen.EXPECT().Generate().Return("generated-secret", nil)
		cfg, err := ProvideConfig(noopLoggerFx{}, keyGen)
		require.NoError(t, err)
		require.NotNil(t, cfg)
		require.Equal(t, "generated-secret", cfg.Secret.Key)
	})

	t.Run("error when key generation fails", func(t *testing.T) {
		setRequiredSettingEnv(t)
		t.Setenv("SECRET_KEY", "")
		t.Setenv("SETTING_SECRET_KEY", "")

		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		keyGen := mocks.NewMockGenerator(ctrl)
		keyGen.EXPECT().Generate().Return("", errors.New("keygen failed"))
		_, err := ProvideConfig(noopLoggerFx{}, keyGen)
		require.Error(t, err)
	})
}

func TestInitObservabilityRegistersOnStopHook(t *testing.T) {
	lc := &fakeLifecycle{}
	cfg := &config.Config{
		General: config.GeneralConfig{
			Env: "test",
		},
		Observability: config.ObservabilityConfig{
			OtelExporterOTLPEndpoint: "localhost:4318",
			OtelServiceName:          "aion-test",
			OtelServiceVersion:       "v1",
			OtelExporterInsecure:     true,
			OtelExporterTimeout:      "1s",
		},
	}

	InitObservability(lc, cfg, noopLoggerFx{})

	require.Len(t, lc.hooks, 1)
	require.NotNil(t, lc.hooks[0].OnStop)
	require.NoError(t, lc.hooks[0].OnStop(t.Context()))
}

func TestProvideHTTPClientDefaultTimeout(t *testing.T) {
	cfg := &config.Config{
		AionChat: config.AionChatConfig{
			Timeout: 0,
		},
	}

	client := ProvideHTTPClient(cfg)
	require.NotNil(t, client)
}

func TestProvideHTTPClientCustomTimeout(t *testing.T) {
	cfg := &config.Config{
		AionChat: config.AionChatConfig{
			Timeout: 250 * time.Millisecond,
		},
	}

	client := ProvideHTTPClient(cfg)
	require.NotNil(t, client)
}

func TestProvideCacheReturnsErrorWhenRedisUnavailable(t *testing.T) {
	lc := &fakeLifecycle{}
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Addr:           "127.0.0.1:1",
			AuthDB:         0,
			CategoryDB:     1,
			TagDB:          2,
			RecordDB:       3,
			UserDB:         4,
			ChatDB:         5,
			PoolSize:       1,
			ConnectTimeout: 10 * time.Millisecond,
		},
	}

	out, err := ProvideCache(lc, cfg, noopLoggerFx{})
	require.Error(t, err)
	require.Equal(t, CacheOut{}, out)
	require.Empty(t, lc.hooks)
}

func TestProvideDatabaseReturnsErrorWhenPostgresUnavailable(t *testing.T) {
	lc := &fakeLifecycle{}
	cfg := &config.Config{
		DB: config.DBConfig{
			Type:            "postgres",
			Host:            "127.0.0.1",
			Port:            "1",
			Name:            "aion",
			User:            "aion",
			Password:        "aion",
			SSLMode:         "disable",
			TimeZone:        "UTC",
			MaxOpenConns:    1,
			MaxIdleConns:    1,
			MaxRetries:      1,
			RetryInterval:   10 * time.Millisecond,
			ConnMaxLifetime: time.Second,
		},
	}

	dbConn, err := ProvideDatabase(lc, cfg, noopLoggerFx{})
	require.Error(t, err)
	require.Nil(t, dbConn)
	require.Empty(t, lc.hooks)
}
