package category

import (
	"context"
	"errors"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

type CategoryRetriever interface {
	GetCategoryByID(ctx context.Context, id uint64) (domain.Category, error)
	GetCategoryByName(ctx context.Context, name string) (domain.Category, error)
	GetAllCategories(ctx context.Context) ([]domain.Category, error)
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, categoryID uint64) (domain.Category, error) {
	if categoryID == 0 {
		s.Logger.Errorw(constants.CategoryIDIsRequired, constants.CategoryID, categoryID)
		return domain.Category{}, errors.New(constants.CategoryIDIsRequired)
	}

	category, err := s.CategoryRepository.GetCategoryByID(ctx, categoryID)
	if err != nil {
		s.Logger.Errorw(constants.FailedToGetCategoryByID, constants.CategoryID, categoryID, constants.Error, err)
		return domain.Category{}, err
	}

	return category, nil
}

func (s *CategoryService) GetCategoryByName(ctx context.Context, name string) (domain.Category, error) {
	if name == "" {
		s.Logger.Errorw(constants.CategoryNameIsRequired, constants.CategoryName, name)
		return domain.Category{}, errors.New(constants.CategoryNameIsRequired)
	}

	category, err := s.CategoryRepository.GetCategoryByName(ctx, name)
	if err != nil {
		s.Logger.Errorw(constants.FailedToGetCategoryByName, constants.CategoryName, name, constants.Error, err)
		return domain.Category{}, err
	}

	return category, nil
}

func (s *CategoryService) GetAllCategories(ctx context.Context) ([]domain.Category, error) {
	categories, err := s.CategoryRepository.GetAllCategories(ctx)
	if err != nil {
		s.Logger.Errorw(constants.FailedToGetAllCategories, constants.Error, err)
		return nil, err
	}

	return categories, nil
}
