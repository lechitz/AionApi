package category

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

// Deleter defines the contract for deleting a category with a soft-delete mechanism.
type Deleter interface {
	SoftDeleteCategory(ctx context.Context, category domain.Category) error
}

// SoftDeleteCategory performs a soft delete operation on a category, marking it as inactive instead of permanently removing it from the database.
func (s *Service) SoftDeleteCategory(ctx context.Context, category domain.Category) error {
	categoryDB, err := s.Repository.GetCategoryByID(ctx, category)
	if err != nil {
		s.Logger.Errorw(
			constants.FailedToGetCategoryByID,
			constants.CategoryID,
			category.ID,
			constants.Error,
			err,
		)
		return errors.New(constants.FailedToGetCategoryByID)
	}

	if err := s.Repository.SoftDeleteCategory(ctx, categoryDB); err != nil {
		s.Logger.Errorw(
			constants.FailedToSoftDeleteCategory,
			constants.CategoryID,
			category.ID,
			constants.Error,
			err,
		)
		return err
	}

	s.Logger.Infow(constants.SuccessfullySoftDeletedCategory, constants.CategoryID, category.ID)
	return nil
}
