package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/aion-api/internal/category/core/domain"
	"github.com/lechitz/aion-api/internal/category/core/usecase"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetByID_InvalidCategoryID(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	categoryDB, err := suite.CategoryService.GetByID(suite.Ctx, 0, 2)

	require.Error(t, err)
	require.Equal(t, usecase.CategoryIDIsRequired, err.Error())
	require.Equal(t, domain.Category{}, categoryDB)
}

func TestGetByID_ErrorToGetByID(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := domain.Category{
		ID:     1,
		UserID: 3,
		Name:   "Work",
	}

	suite.CategoryCache.EXPECT().
		GetCategory(gomock.Any(), category.ID, category.UserID).
		Return(domain.Category{}, errors.New("cache miss"))

	suite.CategoryRepository.EXPECT().
		GetByID(gomock.Any(), category.ID, category.UserID).
		Return(domain.Category{}, errors.New(usecase.FailedToGetCategoryByID))

	_, err := suite.CategoryService.GetByID(suite.Ctx, category.ID, category.UserID)

	require.Error(t, err)
	require.Equal(t, usecase.FailedToGetCategoryByID, err.Error())
}

func TestGetByID_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := domain.Category{
		ID:     1,
		UserID: 3,
		Name:   "Work",
	}

	suite.CategoryCache.EXPECT().
		GetCategory(gomock.Any(), category.ID, category.UserID).
		Return(domain.Category{}, errors.New("cache miss"))

	suite.CategoryRepository.EXPECT().
		GetByID(gomock.Any(), category.ID, category.UserID).
		Return(category, nil)

	suite.CategoryCache.EXPECT().
		SaveCategory(gomock.Any(), category, gomock.Any()).
		Return(nil)

	categoryDB, err := suite.CategoryService.GetByID(suite.Ctx, category.ID, category.UserID)

	require.NoError(t, err)
	require.Equal(t, category, categoryDB)
}
