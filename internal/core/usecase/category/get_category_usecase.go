package category

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/core/domain/entity"
	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

// Retriever defines methods for retrieving category information from a data source.
// It provides lookup by ID, by name, or retrieves all categories for a specific user.
type Retriever interface {
	GetCategoryByID(ctx context.Context, category entity.Category) (entity.Category, error)
	GetCategoryByName(ctx context.Context, category entity.Category) (entity.Category, error)
	GetAllCategories(ctx context.Context, userID uint64) ([]entity.Category, error)
}

// GetCategoryByID retrieves a category by its ID from the database and returns it.
// Returns an error if the ID is invalid or if the retrieval fails.
func (s *Service) GetCategoryByID(ctx context.Context, category entity.Category) (entity.Category, error) {
	if category.ID == 0 {
		s.Logger.Errorw(constants.CategoryIDIsRequired, def.CtxCategoryID, category.ID)
		return entity.Category{}, errors.New(constants.CategoryIDIsRequired)
	}

	categoryDB, err := s.Repository.GetCategoryByID(ctx, category)
	if err != nil {
		s.Logger.Errorw(constants.FailedToGetCategoryByID, def.CtxCategoryID, categoryDB.ID, def.Error, err)
		return entity.Category{}, errors.New(constants.FailedToGetCategoryByID)
	}

	return categoryDB, nil
}

// GetCategoryByName retrieves a category by its name from the database and returns it.
// Returns an error if the name is empty or if the retrieval fails.
func (s *Service) GetCategoryByName(ctx context.Context, category entity.Category) (entity.Category, error) {
	if category.Name == "" {
		s.Logger.Errorw(constants.CategoryNameIsRequired, def.CtxCategoryName, category.Name)
		return entity.Category{}, errors.New(constants.CategoryNameIsRequired)
	}

	categoryDB, err := s.Repository.GetCategoryByName(ctx, category)
	if err != nil {
		s.Logger.Errorw(constants.FailedToGetCategoryByName, def.CtxCategoryName, category.Name, def.Error, err)
		return entity.Category{}, err
	}

	return categoryDB, nil
}

// GetAllCategories retrieves all categories associated with a specific user ID using the repository. Returns a list of categories or an error in case of failure.
func (s *Service) GetAllCategories(ctx context.Context, userID uint64) ([]entity.Category, error) {
	categories, err := s.Repository.GetAllCategories(ctx, userID)
	if err != nil {
		s.Logger.Errorw(constants.FailedToGetAllCategories, def.Error, err)
		return nil, err
	}

	return categories, nil
}
