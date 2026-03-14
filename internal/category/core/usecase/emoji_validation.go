package usecase

import (
	"regexp"
	"strings"
)

var iconKeyPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9/_-]*\.svg$`)

func normalizeIconKey(icon *string) string {
	if icon == nil {
		return ""
	}
	return strings.TrimSpace(*icon)
}

func isValidIconKey(value string) bool {
	if value == "" {
		return false
	}
	if len(value) > 120 {
		return false
	}
	return iconKeyPattern.MatchString(value)
}
