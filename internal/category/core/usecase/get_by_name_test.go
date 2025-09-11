package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/category/core/domain"
	"github.com/lechitz/AionApi/internal/category/core/usecase"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/lechitz/AionApi/tests/testdata"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetCategoryByName_InvalidCategoryName(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	categoryDB, err := suite.CategoryService.GetByName(suite.Ctx, "", category.UserID)

	require.Error(t, err)
	require.Equal(t, usecase.CategoryNameIsRequired, err.Error())
	require.Equal(t, domain.Category{}, categoryDB)
}

func TestGetCategoryByName_ErrorToGetCategoryByName(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

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

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetByName(gomock.Any(), category.Name, category.UserID).
		Return(category, nil)

	categoryDB, err := suite.CategoryService.GetByName(suite.Ctx, category.Name, category.UserID)

	require.NoError(t, err)
	require.Equal(t, category, categoryDB)
}
