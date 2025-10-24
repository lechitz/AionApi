// Package observability contains utilities for configuring OpenTelemetry exporters.
package observability

import (
	"net/url"
	"strings"
)

// ParseHeaders parses a CSV string into a map for OTEL headers.
func ParseHeaders(headersStr string) map[string]string {
	result := make(map[string]string)
	if headersStr == "" {
		return result
	}
	pairs := strings.Split(headersStr, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			result[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return result
}

// NormalizeEndpoint accepts an OTLP endpoint value which may be either
// host:port or a full URL (http(s)://host:port[/path]). It ensures the
// value is a valid URL string. If the provided value has no scheme, it
// temporarily prefixes http:// for parsing/validation and returns the
// resulting string (with scheme). This avoids blind concatenation that
// can lead to encoded values like http:%2F%2F... when passed through
// templates/env expansion.
func NormalizeEndpoint(endpoint string) (string, error) {
	if strings.TrimSpace(endpoint) == "" {
		return "", nil
	}
	e := strings.TrimSpace(endpoint)
	if !strings.HasPrefix(e, "http://") && !strings.HasPrefix(e, "https://") {
		// Temporarily add a scheme to allow url.Parse to validate and normalize.
		e = "http://" + e
	}
	u, err := url.Parse(e)
	if err != nil {
		return "", err
	}
	// Return the normalized string form (includes scheme). Callers may
	// still choose to strip the scheme if their exporter expects host:port.
	return u.String(), nil
}
