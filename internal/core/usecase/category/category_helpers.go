// Package category provides use-case implementations for managing categories.
package category

import (
	"errors"

	"github.com/lechitz/AionApi/internal/core/domain"

	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

// validateCreateCategoryRequired validates required fields for creating a category and enforces constraints like name presence and field length limits.
func (s *Service) validateCreateCategoryRequired(category domain.Category) error {
	if category.Name == "" {
		return errors.New(constants.CategoryNameIsRequired)
	}

	if category.Description != "" && len(category.Description) > 200 {
		return errors.New(constants.CategoryDescriptionIsTooLong)
	}

	if category.Color != "" && len(category.Color) > 7 {
		return errors.New(constants.CategoryColorIsTooLong)
	}

	if category.Icon != "" && len(category.Icon) > 50 {
		return errors.New(constants.CategoryIconIsTooLong)
	}

	return nil
}
