package handlerhelpers

import "github.com/lechitz/AionApi/internal/shared/sharederrors"

// CheckRequiredFields checks if all required fields are present in the given map.
func CheckRequiredFields(fields map[string]string) error {
	var missing []string
	for name, value := range fields {
		if value == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return sharederrors.MissingFields(missing...)
	}
	return nil
}

// AtLeastOneField checks if at least one field is present in the given map.
func AtLeastOneField(fields map[string]interface{}) error {
	var names []string
	for name, value := range fields {
		switch v := value.(type) {
		case *string:
			if v != nil && *v != "" {
				return nil
			}
		case string:
			if v != "" {
				return nil
			}
		}
		names = append(names, name)
	}
	return sharederrors.AtLeastOneFieldRequired(names...)
}

// SafeStringPtr safely dereferences a *string, returning an empty string if nil.
func SafeStringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
