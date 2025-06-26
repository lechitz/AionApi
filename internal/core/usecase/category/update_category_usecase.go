package category

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain/entity"
	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

// Updater defines an interface for updating an existing category in the system.
type Updater interface {
	UpdateCategory(ctx context.Context, category entity.Category) (entity.Category, error)
}

// UpdateCategory updates an existing category in the system with provided fields and logs the operation outcome. Returns the updated category or an error.
func (s *Service) UpdateCategory(ctx context.Context, category entity.Category) (entity.Category, error) {
	fieldsToUpdate := extractUpdateFields(category)

	updatedCategory, err := s.Repository.UpdateCategory(ctx, category.ID, category.UserID, fieldsToUpdate)
	if err != nil {
		s.Logger.Errorw(constants.FailedToUpdateCategory, def.CtxCategoryID, category.ID, def.Error, err)
		return entity.Category{}, err
	}

	s.Logger.Infow(constants.SuccessfullyUpdatedCategory, def.CtxCategoryID, updatedCategory.ID)

	return updatedCategory, nil
}

// extractUpdateFields constructs a map of non-empty category fields for updating.
func extractUpdateFields(category entity.Category) map[string]interface{} {
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
