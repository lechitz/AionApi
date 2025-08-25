package category_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
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
	require.Equal(t, constants.CategoryIDIsRequired, err.Error())
	require.Equal(t, domain.Category{}, categoryDB)
}

func TestGetByID_ErrorToGetByID(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	id := testdata.PerfectCategory.ID

	suite.CategoryRepository.EXPECT().
		GetByID(gomock.Any(), id).
		Return(domain.Category{}, errors.New(constants.FailedToGetCategoryByID))

	_, err := suite.CategoryService.GetByID(suite.Ctx, id)

	require.Error(t, err)
	require.Equal(t, constants.FailedToGetCategoryByID, err.Error())
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
