// Package main boots the outbox publisher command-line application.
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	envBootstrapStartTimeout = "BOOTSTRAP_START_TIMEOUT"
	envBootstrapStopTimeout  = "BOOTSTRAP_STOP_TIMEOUT"
)

const (
	defaultStartTimeout = 20 * time.Second
	defaultStopTimeout  = 20 * time.Second
)

var ErrInvalidBootstrapConfig = errors.New("invalid bootstrap config")

type bootstrapConfig struct {
	StartTimeout time.Duration
	StopTimeout  time.Duration
}

func loadBootstrapConfig(getenv func(string) string) (bootstrapConfig, error) {
	startTimeout, err := readDurationEnv(getenv, envBootstrapStartTimeout, defaultStartTimeout)
	if err != nil {
		return bootstrapConfig{}, err
	}
	stopTimeout, err := readDurationEnv(getenv, envBootstrapStopTimeout, defaultStopTimeout)
	if err != nil {
		return bootstrapConfig{}, err
	}

	return bootstrapConfig{
		StartTimeout: startTimeout,
		StopTimeout:  stopTimeout,
	}, nil
}

func readDurationEnv(getenv func(string) string, key string, defaultValue time.Duration) (time.Duration, error) {
	raw := strings.TrimSpace(getenv(key))
	if raw == "" {
		return defaultValue, nil
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("%w: invalid %s: %w", ErrInvalidBootstrapConfig, key, err)
	}
	if value <= 0 {
		return 0, fmt.Errorf("%w: invalid %s: must be greater than 0", ErrInvalidBootstrapConfig, key)
	}

	return value, nil
}
