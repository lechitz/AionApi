package category

import (
	"context"
	"errors"
	"fmt"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

// CreateCategory creates a new category in the database after validating inputs and ensuring uniqueness by name. Returns the created category or an error.
func (s *Service) CreateCategory(ctx context.Context, category domain.Category) (domain.Category, error) {
	if err := s.validateCreateCategoryRequired(category); err != nil {
		s.Logger.Errorw(constants.ErrToValidateCategory, commonkeys.Error, err.Error())
		return domain.Category{}, err
	}

	existingCategory, err := s.CategoryRepository.GetCategoryByName(ctx, category)
	if err == nil && existingCategory.Name != "" {
		s.Logger.Errorw(constants.CategoryAlreadyExists, commonkeys.CategoryName, category.Name)
		return domain.Category{}, errors.New(constants.CategoryAlreadyExists)
	}

	createdCategory, err := s.CategoryRepository.CreateCategory(ctx, category)
	if err != nil {
		s.Logger.Errorw(constants.FailedToCreateCategory, commonkeys.Category, category, commonkeys.Error, err)
		return domain.Category{}, fmt.Errorf("%s: %w", constants.FailedToCreateCategory, err)
	}

	s.Logger.Infow(fmt.Sprintf(constants.SuccessfullyCreatedCategory, category.Name))

	return createdCategory, nil
}
