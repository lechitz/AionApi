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

func TestGetByID_InvalidCategoryID(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	categoryDB, err := suite.CategoryService.GetByID(suite.Ctx, 1, 2)

	require.Error(t, err)
	require.Equal(t, usecase.CategoryIDIsRequired, err.Error())
	require.Equal(t, domain.Category{}, categoryDB)
}

func TestGetByID_ErrorToGetByID(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

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

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetByID(gomock.Any(), category.ID, category.UserID).
		Return(category, nil)

	categoryDB, err := suite.CategoryService.GetByID(suite.Ctx, category.ID, category.UserID)

	require.NoError(t, err)
	require.Equal(t, category, categoryDB)
}
