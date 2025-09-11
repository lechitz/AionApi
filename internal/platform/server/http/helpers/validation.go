package helpers

import (
	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/sharederrors"
)

//TODO: CheckRequiredFields pode passar pra validação no DTO.

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
