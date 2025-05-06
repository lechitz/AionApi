package category_test

import (
	"errors"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/lechitz/AionApi/tests/testdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateCategory_ErrorToUpdateCategory(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
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

	assert.Error(t, err)
	assert.Equal(t, domain.Category{}, categoryDB)
	assert.Equal(t, constants.FailedToUpdateCategory, err.Error())
}

func TestUpdateCategory_Success(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
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

	assert.NoError(t, err)

	assert.Equal(t, category.ID, categoryDB.ID)
	assert.Equal(t, category.UserID, categoryDB.UserID)
	assert.Equal(t, category.Name, categoryDB.Name)
	assert.Equal(t, category.Description, categoryDB.Description)
	assert.Equal(t, category.Color, categoryDB.Color)
	assert.Equal(t, category.Icon, categoryDB.Icon)
}
