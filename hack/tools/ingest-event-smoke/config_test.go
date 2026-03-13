package main

import "testing"

func TestLoadConfigDefaults(t *testing.T) {
	t.Parallel()

	cfg := loadConfig(nil)
	if cfg.ingestHost != defaultIngestHost || cfg.kafkaBroker != defaultKafkaBroker || cfg.topic != defaultIngestTopic {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}
