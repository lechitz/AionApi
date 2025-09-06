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

func TestGetByID_InvalidCategoryID(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	categoryDB, err := suite.CategoryService.GetByID(suite.Ctx, 0)

	require.Error(t, err)
	require.Equal(t, usecase.CategoryIDIsRequired, err.Error())
	require.Equal(t, domain.Category{}, categoryDB)
}

func TestGetByID_ErrorToGetByID(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	id := testdata.PerfectCategory.ID

	suite.CategoryRepository.EXPECT().
		GetByID(gomock.Any(), id).
		Return(domain.Category{}, errors.New(usecase.FailedToGetCategoryByID))

	_, err := suite.CategoryService.GetByID(suite.Ctx, id)

	require.Error(t, err)
	require.Equal(t, usecase.FailedToGetCategoryByID, err.Error())
}

func TestGetByID_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	exp := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetByID(gomock.Any(), exp.ID).
		Return(exp, nil)

	categoryDB, err := suite.CategoryService.GetByID(suite.Ctx, exp.ID)

	require.NoError(t, err)
	require.Equal(t, exp, categoryDB)
}
