package category

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

// Updater defines an interface for updating an existing category in the system.
type Updater interface {
	UpdateCategory(ctx context.Context, category domain.Category) (domain.Category, error)
}

// UpdateCategory updates an existing category in the system with provided fields and logs the operation outcome. Returns the updated category or an error.
func (s *Service) UpdateCategory(
	ctx context.Context,
	category domain.Category,
) (domain.Category, error) {
	fieldsToUpdate := extractUpdateFields(category)

	updatedCategory, err := s.Repository.UpdateCategory(
		ctx,
		category.ID,
		category.UserID,
		fieldsToUpdate,
	)
	if err != nil {
		s.Logger.Errorw(
			constants.FailedToUpdateCategory,
			constants.CategoryID,
			category.ID,
			constants.Error,
			err,
		)
		return domain.Category{}, err
	}

	s.Logger.Infow(constants.SuccessfullyUpdatedCategory, constants.CategoryID, updatedCategory.ID)
	return updatedCategory, nil
}

// extractUpdateFields constructs a map of non-empty category fields for updating.
func extractUpdateFields(category domain.Category) map[string]interface{} {
	updateFields := make(map[string]interface{})

	if category.Name != "" {
		updateFields[constants.CategoryName] = category.Name
	}
	if category.Description != "" {
		updateFields[constants.CategoryDescription] = category.Description
	}
	if category.Color != "" {
		updateFields[constants.CategoryColor] = category.Color
	}
	if category.Icon != "" {
		updateFields[constants.CategoryIcon] = category.Icon
	}

	return updateFields
}
