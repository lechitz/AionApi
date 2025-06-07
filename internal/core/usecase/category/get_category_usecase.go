package category

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

type CategoryRetriever interface {
	GetCategoryByID(ctx context.Context, category domain.Category) (domain.Category, error)
	GetCategoryByName(ctx context.Context, category domain.Category) (domain.Category, error)
	GetAllCategories(ctx context.Context, userID uint64) ([]domain.Category, error)
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, category domain.Category) (domain.Category, error) {
	if category.ID == 0 {
		s.Logger.Errorw(constants.CategoryIDIsRequired, constants.CategoryID, category.ID)
		return domain.Category{}, errors.New(constants.CategoryIDIsRequired)
	}

	categoryDB, err := s.CategoryRepository.GetCategoryByID(ctx, category)
	if err != nil {
		s.Logger.Errorw(constants.FailedToGetCategoryByID, constants.CategoryID, categoryDB.ID, constants.Error, err)
		return domain.Category{}, err
	}

	return categoryDB, nil
}

func (s *CategoryService) GetCategoryByName(ctx context.Context, category domain.Category) (domain.Category, error) {
	if category.Name == "" {
		s.Logger.Errorw(constants.CategoryNameIsRequired, constants.CategoryName, category.Name)
		return domain.Category{}, errors.New(constants.CategoryNameIsRequired)
	}

	categoryDB, err := s.CategoryRepository.GetCategoryByName(ctx, category)
	if err != nil {
		s.Logger.Errorw(
			constants.FailedToGetCategoryByName,
			constants.CategoryName,
			category.Name,
			constants.Error,
			err,
		)
		return domain.Category{}, err
	}

	return categoryDB, nil
}

func (s *CategoryService) GetAllCategories(ctx context.Context, userID uint64) ([]domain.Category, error) {
	categories, err := s.CategoryRepository.GetAllCategories(ctx, userID)
	if err != nil {
		s.Logger.Errorw(constants.FailedToGetAllCategories, constants.Error, err)
		return nil, err
	}

	return categories, nil
}
