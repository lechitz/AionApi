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

func TestGetCategoryByID_InvalidCategoryID(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	categoryID := uint64(0)

	categoryDB, err := suite.CategoryService.GetCategoryByID(suite.Ctx, categoryID)

	assert.Error(t, err)
	assert.Equal(t, domain.Category{}, categoryDB)
	assert.Equal(t, constants.CategoryIDIsRequired, err.Error())
}

func TestGetCategoryByID_ErrorToGetCategoryByID(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	categoryID := testdata.TestPerfectCategory.ID

	suite.CategoryRepository.EXPECT().
		GetCategoryByID(suite.Ctx, categoryID).
		Return(domain.Category{}, errors.New(constants.FailedToGetCategoryByID))

	_, err := suite.CategoryService.GetCategoryByID(suite.Ctx, categoryID)

	assert.Error(t, err)
	assert.Equal(t, constants.FailedToGetCategoryByID, err.Error())
}

func TestGetCategoryByID_ErrorToCreateCategory(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	categoryID := testdata.TestPerfectCategory.ID

	suite.CategoryRepository.EXPECT().
		GetCategoryByID(suite.Ctx, categoryID).
		Return(domain.Category{}, errors.New(constants.FailedToCreateCategory))

	categoryDB, err := suite.CategoryService.GetCategoryByID(suite.Ctx, categoryID)

	assert.Error(t, err)
	assert.Equal(t, domain.Category{}, categoryDB)
	assert.Equal(t, constants.FailedToCreateCategory, err.Error())
}

func TestGetCategoryByID_Success(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.CategoryRepository.EXPECT().
		GetCategoryByID(suite.Ctx, testdata.TestPerfectCategory.ID).
		Return(testdata.TestPerfectCategory, nil)

	categoryDB, err := suite.CategoryService.GetCategoryByID(suite.Ctx, testdata.TestPerfectCategory.ID)

	assert.NoError(t, err)
	assert.NotNil(t, categoryDB)
	assert.Equal(t, testdata.TestPerfectCategory.ID, categoryDB.ID)
}

func TestGetCategoryByName_InvalidCategoryName(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	categoryName := ""

	categoryDB, err := suite.CategoryService.GetCategoryByName(suite.Ctx, categoryName)

	assert.Error(t, err)
	assert.Equal(t, domain.Category{}, categoryDB)
	assert.Equal(t, constants.CategoryNameIsRequired, err.Error())
}

func TestGetCategoryByName_ErrorToGetCategoryByName(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	categoryName := testdata.TestPerfectCategory.Name

	suite.CategoryRepository.EXPECT().
		GetCategoryByName(suite.Ctx, categoryName).
		Return(domain.Category{}, errors.New(constants.FailedToGetCategoryByName))

	categoryDB, err := suite.CategoryService.GetCategoryByName(suite.Ctx, categoryName)

	assert.Error(t, err)
	assert.Equal(t, domain.Category{}, categoryDB)
	assert.Equal(t, constants.FailedToGetCategoryByName, err.Error())
}

func TestGetCategoryByName_Success(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	CategoryName := testdata.TestPerfectCategory.Name

	suite.CategoryRepository.EXPECT().
		GetCategoryByName(suite.Ctx, CategoryName).
		Return(testdata.TestPerfectCategory, nil)

	categoryDB, err := suite.CategoryService.GetCategoryByName(suite.Ctx, CategoryName)

	assert.NoError(t, err)
	assert.NotNil(t, domain.Category{}, categoryDB)
}

func TestGetAllCategories_ErrorToGetAllCategories(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.CategoryRepository.EXPECT().
		GetAllCategories(suite.Ctx).
		Return(nil, errors.New(constants.FailedToGetAllCategories))

	categories, err := suite.CategoryService.GetAllCategories(suite.Ctx)

	assert.Error(t, err)
	assert.Nil(t, categories)
	assert.Equal(t, constants.FailedToGetAllCategories, err.Error())
}

func TestGetAllCategories_Success(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.CategoryRepository.EXPECT().
		GetAllCategories(suite.Ctx).
		Return([]domain.Category{testdata.TestPerfectCategory}, nil)

	categories, err := suite.CategoryService.GetAllCategories(suite.Ctx)

	assert.NoError(t, err)
	assert.NotNil(t, categories)
	assert.Len(t, categories, 1)
	assert.Equal(t, testdata.TestPerfectCategory, categories[0])
}
