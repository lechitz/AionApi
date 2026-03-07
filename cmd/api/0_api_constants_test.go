package main

import "testing"

func TestBootstrapDefaultsArePositive(t *testing.T) {
	if defaultStartTimeout <= 0 {
		t.Fatalf("defaultStartTimeout must be positive, got %v", defaultStartTimeout)
	}
	if defaultStopTimeout <= 0 {
		t.Fatalf("defaultStopTimeout must be positive, got %v", defaultStopTimeout)
	}
}

func TestBootstrapEnvKeysAreDefined(t *testing.T) {
	if envBootstrapStartTimeout == "" {
		t.Fatal("envBootstrapStartTimeout must not be empty")
	}
	if envBootstrapStopTimeout == "" {
		t.Fatal("envBootstrapStopTimeout must not be empty")
	}
}
