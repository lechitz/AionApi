package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/category/domain"
	"github.com/lechitz/AionApi/internal/core/category/usecase"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/lechitz/AionApi/tests/testdata"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetCategoryByName_InvalidCategoryName(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	categoryDB, err := suite.CategoryService.GetByName(suite.Ctx, "")

	require.Error(t, err)
	require.Equal(t, usecase.CategoryNameIsRequired, err.Error())
	require.Equal(t, domain.Category{}, categoryDB)
}

func TestGetCategoryByName_ErrorToGetCategoryByName(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	name := testdata.PerfectCategory.Name

	suite.CategoryRepository.EXPECT().
		GetByName(gomock.Any(), name).
		Return(domain.Category{}, errors.New(usecase.FailedToGetCategoryByName))

	categoryDB, err := suite.CategoryService.GetByName(suite.Ctx, name)

	require.Error(t, err)
	require.Equal(t, usecase.FailedToGetCategoryByName, err.Error())
	require.Equal(t, domain.Category{}, categoryDB)
}

func TestGetCategoryByName_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	exp := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetByName(gomock.Any(), exp.Name).
		Return(exp, nil)

	categoryDB, err := suite.CategoryService.GetByName(suite.Ctx, exp.Name)

	require.NoError(t, err)
	require.Equal(t, exp, categoryDB)
}
