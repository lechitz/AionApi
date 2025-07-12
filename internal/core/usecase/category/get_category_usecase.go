package category

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

// GetCategoryByID retrieves a category by its ID from the database and returns it.
func (s *Service) GetCategoryByID(ctx context.Context, category domain.Category) (domain.Category, error) {
	if category.ID == 0 {
		s.Logger.Errorw(constants.CategoryIDIsRequired, commonkeys.CategoryID, strconv.FormatUint(category.ID, 10))
		return domain.Category{}, errors.New(constants.CategoryIDIsRequired)
	}

	categoryDB, err := s.CategoryRepository.GetCategoryByID(ctx, category)
	if err != nil {
		s.Logger.Errorw(constants.FailedToGetCategoryByID, commonkeys.CategoryID, strconv.FormatUint(category.ID, 10), commonkeys.Error, err.Error())
		return domain.Category{}, errors.New(constants.FailedToGetCategoryByID)
	}

	return categoryDB, nil
}

// GetCategoryByName retrieves a category by its name from the database and returns it.
func (s *Service) GetCategoryByName(ctx context.Context, category domain.Category) (domain.Category, error) {
	if category.Name == "" {
		s.Logger.Errorw(constants.CategoryNameIsRequired, commonkeys.CategoryName, category.Name)
		return domain.Category{}, errors.New(constants.CategoryNameIsRequired)
	}

	categoryDB, err := s.CategoryRepository.GetCategoryByName(ctx, category)
	if err != nil {
		s.Logger.Errorw(constants.FailedToGetCategoryByName, commonkeys.CategoryName, category.Name, commonkeys.Error, err)
		return domain.Category{}, err
	}

	return categoryDB, nil
}

// GetAllCategories retrieves all categories associated with a specific user ID using the repository. Returns a list of categories or an error in case of failure.
func (s *Service) GetAllCategories(ctx context.Context, userID uint64) ([]domain.Category, error) {
	categories, err := s.CategoryRepository.GetAllCategories(ctx, userID)
	if err != nil {
		s.Logger.Errorw(constants.FailedToGetAllCategories, commonkeys.Error, err)
		return nil, err
	}

	return categories, nil
}
