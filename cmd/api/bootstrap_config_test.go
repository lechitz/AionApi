package main

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestLoadBootstrapConfigDefaults(t *testing.T) {
	cfg, err := loadBootstrapConfig(func(string) string { return "" })
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.StartTimeout != defaultStartTimeout {
		t.Fatalf("expected default start timeout %v, got %v", defaultStartTimeout, cfg.StartTimeout)
	}
	if cfg.StopTimeout != defaultStopTimeout {
		t.Fatalf("expected default stop timeout %v, got %v", defaultStopTimeout, cfg.StopTimeout)
	}
}

func TestLoadBootstrapConfigUsesValidEnv(t *testing.T) {
	cfg, err := loadBootstrapConfig(func(key string) string {
		switch key {
		case envBootstrapStartTimeout:
			return " 21s "
		case envBootstrapStopTimeout:
			return "13s"
		default:
			return ""
		}
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.StartTimeout != 21*time.Second {
		t.Fatalf("expected 21s start timeout, got %v", cfg.StartTimeout)
	}
	if cfg.StopTimeout != 13*time.Second {
		t.Fatalf("expected 13s stop timeout, got %v", cfg.StopTimeout)
	}
}

func TestLoadBootstrapConfigInvalidStartTimeout(t *testing.T) {
	_, err := loadBootstrapConfig(func(key string) string {
		if key == envBootstrapStartTimeout {
			return "abc"
		}
		return ""
	})
	assertBootstrapErr(t, err, envBootstrapStartTimeout)
}

func TestLoadBootstrapConfigZeroStopTimeout(t *testing.T) {
	_, err := loadBootstrapConfig(func(key string) string {
		if key == envBootstrapStopTimeout {
			return "0s"
		}
		return ""
	})
	assertBootstrapErr(t, err, envBootstrapStopTimeout)
}

func TestLoadBootstrapConfigNegativeStartTimeout(t *testing.T) {
	_, err := loadBootstrapConfig(func(key string) string {
		if key == envBootstrapStartTimeout {
			return "-5s"
		}
		return ""
	})
	assertBootstrapErr(t, err, envBootstrapStartTimeout)
}

func TestLoadBootstrapConfigDefaultStopTimeoutWhenOnlyStartSet(t *testing.T) {
	cfg, err := loadBootstrapConfig(func(key string) string {
		if key == envBootstrapStartTimeout {
			return "7s"
		}
		return ""
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.StartTimeout != 7*time.Second {
		t.Fatalf("expected start timeout 7s, got %v", cfg.StartTimeout)
	}
	if cfg.StopTimeout != defaultStopTimeout {
		t.Fatalf("expected default stop timeout %v, got %v", defaultStopTimeout, cfg.StopTimeout)
	}
}

func assertBootstrapErr(t *testing.T, err error, key string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), key) {
		t.Fatalf("expected error to contain %s, got %v", key, err)
	}
	if !errors.Is(err, ErrInvalidBootstrapConfig) {
		t.Fatalf("expected ErrInvalidBootstrapConfig, got %v", err)
	}
}
