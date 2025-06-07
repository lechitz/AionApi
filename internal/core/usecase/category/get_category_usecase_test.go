package category_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/lechitz/AionApi/tests/testdata"
	"github.com/stretchr/testify/assert"
)

func TestGetCategoryByID_InvalidCategoryID(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByID(suite.Ctx, category).
		Return(domain.Category{}, errors.New(constants.CategoryIDIsRequired))

	categoryDB, err := suite.CategoryService.GetCategoryByID(suite.Ctx, category)

	assert.Error(t, err)
	assert.Equal(t, domain.Category{}, categoryDB)
	assert.Equal(t, constants.CategoryIDIsRequired, err.Error())
}

func TestGetCategoryByID_ErrorToGetCategoryByID(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByID(suite.Ctx, category).
		Return(domain.Category{}, errors.New(constants.FailedToGetCategoryByID))

	_, err := suite.CategoryService.GetCategoryByID(suite.Ctx, category)

	assert.Error(t, err)
	assert.Equal(t, constants.FailedToGetCategoryByID, err.Error())
}

func TestGetCategoryByID_ErrorToCreateCategory(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByID(suite.Ctx, category).
		Return(domain.Category{}, errors.New(constants.FailedToCreateCategory))

	categoryDB, err := suite.CategoryService.GetCategoryByID(suite.Ctx, category)

	assert.Error(t, err)
	assert.Equal(t, domain.Category{}, categoryDB)
	assert.Equal(t, constants.FailedToCreateCategory, err.Error())
}

func TestGetCategoryByID_Success(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.CategoryRepository.EXPECT().
		GetCategoryByID(suite.Ctx, testdata.PerfectCategory).
		Return(testdata.PerfectCategory, nil)

	categoryDB, err := suite.CategoryService.GetCategoryByID(suite.Ctx, testdata.PerfectCategory)

	assert.NoError(t, err)
	assert.NotNil(t, categoryDB)
	assert.Equal(t, testdata.PerfectCategory.ID, categoryDB.ID)
}

func TestGetCategoryByName_InvalidCategoryName(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByName(suite.Ctx, category).
		Return(domain.Category{}, errors.New(constants.CategoryNameIsRequired))

	categoryDB, err := suite.CategoryService.GetCategoryByName(suite.Ctx, category)

	assert.Error(t, err)
	assert.Equal(t, domain.Category{}, categoryDB)
	assert.Equal(t, constants.CategoryNameIsRequired, err.Error())
}

func TestGetCategoryByName_ErrorToGetCategoryByName(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByName(suite.Ctx, category).
		Return(domain.Category{}, errors.New(constants.FailedToGetCategoryByName))

	categoryDB, err := suite.CategoryService.GetCategoryByName(suite.Ctx, category)

	assert.Error(t, err)
	assert.Equal(t, domain.Category{}, categoryDB)
	assert.Equal(t, constants.FailedToGetCategoryByName, err.Error())
}

func TestGetCategoryByName_Success(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByName(suite.Ctx, category).
		Return(testdata.PerfectCategory, nil)

	categoryDB, err := suite.CategoryService.GetCategoryByName(suite.Ctx, category)

	assert.NoError(t, err)
	assert.NotNil(t, domain.Category{}, categoryDB)
}

func TestGetAllCategories_ErrorToGetAllCategories(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := testdata.TestPerfectUser.ID

	suite.CategoryRepository.EXPECT().
		GetAllCategories(suite.Ctx, userID).
		Return(nil, errors.New(constants.FailedToGetAllCategories))

	categories, err := suite.CategoryService.GetAllCategories(suite.Ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, categories)
	assert.Equal(t, constants.FailedToGetAllCategories, err.Error())
}

func TestGetAllCategories_Success(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := testdata.TestPerfectUser.ID

	suite.CategoryRepository.EXPECT().
		GetAllCategories(suite.Ctx, userID).
		Return([]domain.Category{testdata.PerfectCategory}, nil)

	categories, err := suite.CategoryService.GetAllCategories(suite.Ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, categories)
	assert.Len(t, categories, 1)
	assert.Equal(t, testdata.PerfectCategory, categories[0])
}
