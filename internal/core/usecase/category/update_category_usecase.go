package category

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/commonkeys"

	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

// Updater defines an interface for updating an existing category in the system.
type Updater interface {
	UpdateCategory(ctx context.Context, category domain.Category) (domain.Category, error)
}

// UpdateCategory updates an existing category in the system with provided fields and logs the operation outcome. Returns the updated category or an error.
func (s *Service) UpdateCategory(ctx context.Context, category domain.Category) (domain.Category, error) {
	fieldsToUpdate := extractUpdateFields(category)

	updatedCategory, err := s.CategoryRepository.UpdateCategory(ctx, category.ID, category.UserID, fieldsToUpdate)
	if err != nil {
		s.Logger.Errorw(constants.FailedToUpdateCategory, commonkeys.CategoryID, strconv.FormatUint(category.ID, 10), commonkeys.Error, err)
		return domain.Category{}, err
	}

	s.Logger.Infow(constants.SuccessfullyUpdatedCategory, commonkeys.CategoryID, strconv.FormatUint(updatedCategory.ID, 10))

	return updatedCategory, nil
}

// extractUpdateFields constructs a map of non-empty category fields for updating.
func extractUpdateFields(category domain.Category) map[string]interface{} {
	updateFields := make(map[string]interface{})

	if category.Name != "" {
		updateFields[commonkeys.CategoryName] = category.Name
	}
	if category.Description != "" {
		updateFields[commonkeys.CategoryDescription] = category.Description
	}
	if category.Color != "" {
		updateFields[commonkeys.CategoryColor] = category.Color
	}
	if category.Icon != "" {
		updateFields[commonkeys.CategoryIcon] = category.Icon
	}

	return updateFields
}
