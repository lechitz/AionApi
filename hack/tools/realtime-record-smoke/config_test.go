package main

import (
	"testing"
	"time"
)

func TestLoadConfigDefaults(t *testing.T) {
	t.Parallel()

	cfg := loadConfig(nil)

	if cfg.host != defaultRealtimeSmokeHost {
		t.Fatalf("host = %q, want %q", cfg.host, defaultRealtimeSmokeHost)
	}
	if cfg.username != defaultRealtimeSmokeUser {
		t.Fatalf("username = %q, want %q", cfg.username, defaultRealtimeSmokeUser)
	}
	if cfg.password != defaultRealtimeSmokePassword {
		t.Fatalf("password = %q, want %q", cfg.password, defaultRealtimeSmokePassword)
	}
	if cfg.tagID != defaultRealtimeSmokeTagID {
		t.Fatalf("tagID = %q, want %q", cfg.tagID, defaultRealtimeSmokeTagID)
	}
	if cfg.timeout != defaultRealtimeSmokeTimeout {
		t.Fatalf("timeout = %v, want %v", cfg.timeout, defaultRealtimeSmokeTimeout)
	}
	if cfg.pollInterval != defaultRealtimeSmokePollInterval {
		t.Fatalf("pollInterval = %v, want %v", cfg.pollInterval, defaultRealtimeSmokePollInterval)
	}
}

func TestLoadConfigOverrides(t *testing.T) {
	t.Parallel()

	cfg := loadConfig([]string{
		"--host", "http://example.local",
		"--username", "alice",
		"--password", "secret",
		"--tag-id", "99",
		"--timeout", "45s",
		"--poll-interval", "2s",
	})

	if cfg.host != "http://example.local" {
		t.Fatalf("host = %q", cfg.host)
	}
	if cfg.username != "alice" {
		t.Fatalf("username = %q", cfg.username)
	}
	if cfg.password != "secret" {
		t.Fatalf("password = %q", cfg.password)
	}
	if cfg.tagID != "99" {
		t.Fatalf("tagID = %q", cfg.tagID)
	}
	if cfg.timeout != 45*time.Second {
		t.Fatalf("timeout = %v", cfg.timeout)
	}
	if cfg.pollInterval != 2*time.Second {
		t.Fatalf("pollInterval = %v", cfg.pollInterval)
	}
}
