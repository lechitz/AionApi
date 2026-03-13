package main

import "testing"

func TestLoadConfigDefaults(t *testing.T) {
	t.Parallel()

	cfg := loadConfig(nil)

	if cfg.host != defaultPageSmokeHost {
		t.Fatalf("expected default host %q, got %q", defaultPageSmokeHost, cfg.host)
	}
	if cfg.username != defaultPageSmokeUser {
		t.Fatalf("expected default username %q, got %q", defaultPageSmokeUser, cfg.username)
	}
	if cfg.password != defaultPageSmokePassword {
		t.Fatalf("expected default password %q, got %q", defaultPageSmokePassword, cfg.password)
	}
	if cfg.tagID != defaultPageSmokeTagID {
		t.Fatalf("expected default tag-id %q, got %q", defaultPageSmokeTagID, cfg.tagID)
	}
	if cfg.pageLimit != defaultPageSmokeLimit {
		t.Fatalf("expected default page-limit %d, got %d", defaultPageSmokeLimit, cfg.pageLimit)
	}
	if cfg.timeout != defaultPageSmokeTimeout {
		t.Fatalf("expected default timeout %s, got %s", defaultPageSmokeTimeout, cfg.timeout)
	}
	if cfg.pollInterval != defaultPageSmokePoll {
		t.Fatalf("expected default poll interval %s, got %s", defaultPageSmokePoll, cfg.pollInterval)
	}
}
