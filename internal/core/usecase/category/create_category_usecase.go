package category

import (
	"context"
	"errors"
	"fmt"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

type CategoryCreator interface {
	CreateCategory(ctx context.Context, category domain.Category) (domain.Category, error)
}

func (s *CategoryService) CreateCategory(ctx context.Context, category domain.Category) (domain.Category, error) {
	if err := s.validateCreateCategoryRequired(category); err != nil {
		s.Logger.Errorw(constants.ErrToValidateCategory, constants.Error, err.Error())
		return domain.Category{}, err
	}

	existingCategory, err := s.CategoryRepository.GetCategoryByName(ctx, category.Name)
	if err == nil && existingCategory.Name != "" {
		s.Logger.Errorw(constants.CategoryAlreadyExists, constants.CategoryName, category.Name)
		return domain.Category{}, errors.New(constants.CategoryAlreadyExists)
	}

	createdCategory, err := s.CategoryRepository.CreateCategory(ctx, category)
	if err != nil {
		s.Logger.Errorw(constants.FailedToCreateCategory, constants.Category, category, constants.Error, err)
		return domain.Category{}, errors.New(fmt.Sprintf(constants.FailedToCreateCategory))
	}

	s.Logger.Infow(fmt.Sprintf(constants.SuccessfullyCreatedCategory, category.Name))

	return createdCategory, nil
}
