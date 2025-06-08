package category

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

type CategoryDeleter interface {
	SoftDeleteCategory(ctx context.Context, category domain.Category) error
}

func (s *CategoryService) SoftDeleteCategory(ctx context.Context, category domain.Category) error {
	categoryDB, err := s.CategoryRepository.GetCategoryByID(ctx, category)
	if err != nil {
		s.Logger.Errorw(constants.FailedToGetCategoryByID, constants.CategoryID, category.ID, constants.Error, err)
		return errors.New(constants.FailedToGetCategoryByID)
	}

	if err := s.CategoryRepository.SoftDeleteCategory(ctx, categoryDB); err != nil {
		s.Logger.Errorw(constants.FailedToSoftDeleteCategory, constants.CategoryID, category.ID, constants.Error, err)
		return err
	}

	s.Logger.Infow(constants.SuccessfullySoftDeletedCategory, constants.CategoryID, category.ID)
	return nil
}
