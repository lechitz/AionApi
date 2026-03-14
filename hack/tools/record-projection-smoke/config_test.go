package main

import "testing"

func TestLoadConfigDefaults(t *testing.T) {
	t.Parallel()

	cfg := loadConfig(nil)
	if cfg.host != defaultSmokeHost || cfg.username != defaultSmokeUser || cfg.tagID != defaultSmokeTagID {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}
