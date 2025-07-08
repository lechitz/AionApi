// Package observability contains utilities for configuring OpenTelemetry exporters.
package observability

import (
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
