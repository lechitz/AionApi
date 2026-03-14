package main

import "testing"

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	cfg := loadConfig(nil)
	if cfg.sampleLimit != defaultSampleLimit {
		t.Fatalf("expected sampleLimit=%d, got %d", defaultSampleLimit, cfg.sampleLimit)
	}
	if cfg.envFile != defaultEnvFile {
		t.Fatalf("expected envFile=%q, got %q", defaultEnvFile, cfg.envFile)
	}

	cfg = loadConfig([]string{"-sample-limit", "9", "-env-file", "custom.env"})
	if cfg.sampleLimit != 9 || cfg.envFile != "custom.env" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}
