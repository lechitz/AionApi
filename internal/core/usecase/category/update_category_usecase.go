package category

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

type CategoryUpdater interface {
	UpdateCategory(ctx context.Context, category domain.Category) (domain.Category, error)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, category domain.Category) (domain.Category, error) {
	fieldsToUpdate := extractUpdateFields(category)

	updatedCategory, err := s.CategoryRepository.UpdateCategory(ctx, category.ID, category.UserID, fieldsToUpdate)
	if err != nil {
		s.Logger.Errorw(constants.FailedToUpdateCategory, constants.CategoryID, category.ID, constants.Error, err)
		return domain.Category{}, err
	}

	s.Logger.Infow(constants.SuccessfullyUpdatedCategory, constants.CategoryID, updatedCategory.ID)
	return updatedCategory, nil
}

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
