package category

import (
	"context"
	"errors"
	"fmt"

	"github.com/lechitz/AionApi/internal/core/domain/entity"
	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

// Creator is an interface for creating categories within the system.
// It defines a method to persist a new category with context-aware operations.
type Creator interface {
	CreateCategory(ctx context.Context, category entity.Category) (entity.Category, error)
}

// CreateCategory creates a new category in the database after validating inputs and ensuring uniqueness by name. Returns the created category or an error.
func (s *Service) CreateCategory(ctx context.Context, category entity.Category) (entity.Category, error) {
	if err := s.validateCreateCategoryRequired(category); err != nil {
		s.Logger.Errorw(constants.ErrToValidateCategory, def.Error, err.Error())
		return entity.Category{}, err
	}

	existingCategory, err := s.Repository.GetCategoryByName(ctx, category)
	if err == nil && existingCategory.Name != "" {
		s.Logger.Errorw(constants.CategoryAlreadyExists, def.CtxCategoryName, category.Name)
		return entity.Category{}, errors.New(constants.CategoryAlreadyExists)
	}

	createdCategory, err := s.Repository.CreateCategory(ctx, category)
	if err != nil {
		s.Logger.Errorw(constants.FailedToCreateCategory, def.CtxCategory, category, def.Error, err)
		return entity.Category{}, fmt.Errorf("%s: %w", constants.FailedToCreateCategory, err)
	}

	s.Logger.Infow(fmt.Sprintf(constants.SuccessfullyCreatedCategory, category.Name))

	return createdCategory, nil
}
