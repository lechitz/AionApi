// Package main provides env-loading helpers for the outbox diagnose tool.
package main

import (
	"bufio"
	"os"
	"strings"
)

func loadEnvFile(path string) error {
	if !strings.HasSuffix(path, ".env") && !strings.HasSuffix(path, ".env.dev") && !strings.HasSuffix(path, ".env.local") {
		return os.ErrInvalid
	}

	//nolint:gosec // path comes from explicit operator flag for local diagnostics
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" || os.Getenv(key) != "" {
			continue
		}
		_ = os.Setenv(key, value)
	}

	return scanner.Err()
}
