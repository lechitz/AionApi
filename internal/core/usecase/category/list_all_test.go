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

func TestListAll_ErrorToGetAllCategories(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := testdata.PerfectCategory.UserID

	suite.CategoryRepository.EXPECT().
		ListAll(gomock.Any(), userID).
		Return(nil, errors.New(constants.FailedToGetAllCategories))

	categories, err := suite.CategoryService.ListAll(suite.Ctx, userID)

	require.Error(t, err)
	require.Nil(t, categories)
	require.Equal(t, constants.FailedToGetAllCategories, err.Error())
}

func TestListAll_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := testdata.PerfectCategory.UserID
	exp := []domain.Category{testdata.PerfectCategory}

	suite.CategoryRepository.EXPECT().
		ListAll(gomock.Any(), userID).
		Return(exp, nil)

	categories, err := suite.CategoryService.ListAll(suite.Ctx, userID)

	require.NoError(t, err)
	require.Equal(t, exp, categories)
}
