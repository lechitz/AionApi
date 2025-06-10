package category

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

// Retriever defines methods for retrieving category information from a data source.
// It provides lookup by ID, by name, or retrieves all categories for a specific user.
type Retriever interface {
	GetCategoryByID(ctx context.Context, category domain.Category) (domain.Category, error)
	GetCategoryByName(ctx context.Context, category domain.Category) (domain.Category, error)
	GetAllCategories(ctx context.Context, userID uint64) ([]domain.Category, error)
}

// GetCategoryByID retrieves a category by its ID from the database and returns it.
// Returns an error if the ID is invalid or if the retrieval fails.
func (s *Service) GetCategoryByID(ctx context.Context, category domain.Category) (domain.Category, error) {
	if category.ID == 0 {
		s.Logger.Errorw(constants.CategoryIDIsRequired, constants.CategoryID, category.ID)
		return domain.Category{}, errors.New(constants.CategoryIDIsRequired)
	}

	categoryDB, err := s.Repository.GetCategoryByID(ctx, category)
	if err != nil {
		s.Logger.Errorw(
			constants.FailedToGetCategoryByID,
			constants.CategoryID,
			categoryDB.ID,
			constants.Error,
			err,
		)
		return domain.Category{}, err
	}

	return categoryDB, nil
}

// GetCategoryByName retrieves a category by its name from the database and returns it.
// Returns an error if the name is empty or if the retrieval fails.
func (s *Service) GetCategoryByName(ctx context.Context, category domain.Category) (domain.Category, error) {
	if category.Name == "" {
		s.Logger.Errorw(constants.CategoryNameIsRequired, constants.CategoryName, category.Name)
		return domain.Category{}, errors.New(constants.CategoryNameIsRequired)
	}

	categoryDB, err := s.Repository.GetCategoryByName(ctx, category)
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

// GetAllCategories retrieves all categories associated with a specific user ID using the repository. Returns a list of categories or an error in case of failure.
func (s *Service) GetAllCategories(ctx context.Context, userID uint64) ([]domain.Category, error) {
	categories, err := s.Repository.GetAllCategories(ctx, userID)
	if err != nil {
		s.Logger.Errorw(constants.FailedToGetAllCategories, constants.Error, err)
		return nil, err
	}

	return categories, nil
}
