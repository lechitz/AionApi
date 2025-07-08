package category

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/common"

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
		s.Logger.Errorw(constants.FailedToUpdateCategory, common.CategoryID, strconv.FormatUint(category.ID, 10), common.Error, err)
		return domain.Category{}, err
	}

	s.Logger.Infow(constants.SuccessfullyUpdatedCategory, common.CategoryID, strconv.FormatUint(updatedCategory.ID, 10))

	return updatedCategory, nil
}

// extractUpdateFields constructs a map of non-empty category fields for updating.
func extractUpdateFields(category domain.Category) map[string]interface{} {
	updateFields := make(map[string]interface{})

	if category.Name != "" {
		updateFields[common.CategoryName] = category.Name
	}
	if category.Description != "" {
		updateFields[common.CategoryDescription] = category.Description
	}
	if category.Color != "" {
		updateFields[common.CategoryColor] = category.Color
	}
	if category.Icon != "" {
		updateFields[common.CategoryIcon] = category.Icon
	}

	return updateFields
}
