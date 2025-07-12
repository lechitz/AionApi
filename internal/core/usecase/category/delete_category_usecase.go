package category

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

// SoftDeleteCategory performs a soft delete operation on a category, marking it as inactive instead of permanently removing it from the database.
func (s *Service) SoftDeleteCategory(ctx context.Context, category domain.Category) error {
	categoryDB, err := s.CategoryRepository.GetCategoryByID(ctx, category)
	if err != nil {
		s.Logger.Errorw(constants.FailedToGetCategoryByID, commonkeys.CategoryID, strconv.FormatUint(category.ID, 10), commonkeys.Error, err)
		return errors.New(constants.FailedToGetCategoryByID)
	}

	if err := s.CategoryRepository.SoftDeleteCategory(ctx, categoryDB); err != nil {
		s.Logger.Errorw(constants.FailedToSoftDeleteCategory, commonkeys.CategoryID, strconv.FormatUint(category.ID, 10), commonkeys.Error, err)
		return errors.New(constants.FailedToSoftDeleteCategory)
	}

	s.Logger.Infow(constants.SuccessfullySoftDeletedCategory, commonkeys.CategoryID, strconv.FormatUint(category.ID, 10))

	return nil
}
