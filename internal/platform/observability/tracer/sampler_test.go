//nolint:testpackage // tests need access to unexported sampler/resource helpers.
package tracer

import (
	"testing"

	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/stretchr/testify/require"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
)

func TestConfigureSamplerFromEnv(t *testing.T) {
	t.Setenv("OTEL_TRACES_SAMPLER", "always_on")
	require.NotNil(t, configureSamplerFromEnv())

	t.Setenv("OTEL_TRACES_SAMPLER", "traceidratio")
	t.Setenv("OTEL_TRACES_SAMPLER_ARG", "0.5")
	require.NotNil(t, configureSamplerFromEnv())

	t.Setenv("OTEL_TRACES_SAMPLER", "unknown")
	require.Nil(t, configureSamplerFromEnv())
}

func TestBuildResource(t *testing.T) {
	cfg := &config.Config{
		General: config.GeneralConfig{
			Env: "test",
		},
		Observability: config.ObservabilityConfig{
			OtelServiceName:    "svc",
			OtelServiceVersion: "v1",
		},
	}
	res := buildResource(cfg)
	require.NotNil(t, res)

	attrs := res.Attributes()
	require.NotEmpty(t, attrs)

	foundServiceName := false
	for _, attr := range attrs {
		if string(attr.Key) == string(semconv.ServiceNameKey) && attr.Value.AsString() == "svc" {
			foundServiceName = true
			break
		}
	}
	require.True(t, foundServiceName)
}
