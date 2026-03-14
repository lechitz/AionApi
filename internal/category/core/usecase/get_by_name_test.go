package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/category/core/domain"
	"github.com/lechitz/AionApi/internal/category/core/usecase"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetCategoryByName_InvalidCategoryName(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(3)

	categoryDB, err := suite.CategoryService.GetByName(suite.Ctx, "", userID)

	require.Error(t, err)
	require.Equal(t, usecase.CategoryNameIsRequired, err.Error())
	require.Equal(t, domain.Category{}, categoryDB)
}

func TestGetCategoryByName_ErrorToGetCategoryByName(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := domain.Category{
		UserID: 3,
		Name:   "Work",
	}

	suite.CategoryCache.EXPECT().
		GetCategoryByName(gomock.Any(), category.Name, category.UserID).
		Return(domain.Category{}, errors.New("cache miss"))

	suite.CategoryRepository.EXPECT().
		GetByName(gomock.Any(), category.Name, category.UserID).
		Return(domain.Category{}, errors.New(usecase.FailedToGetCategoryByName))

	categoryDB, err := suite.CategoryService.GetByName(suite.Ctx, category.Name, category.UserID)

	require.Error(t, err)
	require.Equal(t, usecase.FailedToGetCategoryByName, err.Error())
	require.Equal(t, domain.Category{}, categoryDB)
}

func TestGetCategoryByName_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := domain.Category{
		UserID: 3,
		Name:   "Work",
	}

	suite.CategoryCache.EXPECT().
		GetCategoryByName(gomock.Any(), category.Name, category.UserID).
		Return(domain.Category{}, errors.New("cache miss"))

	suite.CategoryRepository.EXPECT().
		GetByName(gomock.Any(), category.Name, category.UserID).
		Return(category, nil)

	suite.CategoryCache.EXPECT().
		SaveCategoryByName(gomock.Any(), category, gomock.Any()).
		Return(nil)

	// SaveCategory is only called if category.ID != 0
	// But in this test, category.ID is 0 (not set), so this expectation should be removed
	// OR we need to set category.ID to a non-zero value

	categoryDB, err := suite.CategoryService.GetByName(suite.Ctx, category.Name, category.UserID)

	require.NoError(t, err)
	require.Equal(t, category, categoryDB)
}
