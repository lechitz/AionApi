package category_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain"

	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/lechitz/AionApi/tests/testdata"
	"github.com/stretchr/testify/require"
)

func TestUpdateCategory_ErrorToUpdateCategory(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	updateFields := map[string]interface{}{
		constants.CategoryName:        testdata.PerfectCategory.Name,
		constants.CategoryDescription: testdata.PerfectCategory.Description,
		constants.CategoryColor:       testdata.PerfectCategory.Color,
		constants.CategoryIcon:        testdata.PerfectCategory.Icon,
	}

	suite.CategoryRepository.EXPECT().
		UpdateCategory(suite.Ctx, testdata.PerfectCategory.ID, testdata.PerfectCategory.UserID, updateFields).
		Return(domain.Category{}, errors.New(constants.FailedToUpdateCategory))

	categoryDB, err := suite.CategoryService.UpdateCategory(suite.Ctx, testdata.PerfectCategory)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, categoryDB)
	require.Equal(t, constants.FailedToUpdateCategory, err.Error())
}

func TestUpdateCategory_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	updateFields := map[string]interface{}{
		constants.CategoryName:        testdata.PerfectCategory.Name,
		constants.CategoryDescription: testdata.PerfectCategory.Description,
		constants.CategoryColor:       testdata.PerfectCategory.Color,
		constants.CategoryIcon:        testdata.PerfectCategory.Icon,
	}

	suite.CategoryRepository.EXPECT().
		UpdateCategory(suite.Ctx, testdata.PerfectCategory.ID, testdata.PerfectCategory.UserID, updateFields).
		Return(testdata.PerfectCategory, nil)

	categoryDB, err := suite.CategoryService.UpdateCategory(suite.Ctx, testdata.PerfectCategory)

	require.NoError(t, err)
	require.Equal(t, category.ID, categoryDB.ID)
	require.Equal(t, category.UserID, categoryDB.UserID)
	require.Equal(t, category.Name, categoryDB.Name)
	require.Equal(t, category.Description, categoryDB.Description)
	require.Equal(t, category.Color, categoryDB.Color)
	require.Equal(t, category.Icon, categoryDB.Icon)
}
