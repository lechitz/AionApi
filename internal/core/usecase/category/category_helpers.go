package category

import (
	"errors"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

func (s *CategoryService) validateCreateCategoryRequired(category domain.Category) error {
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
