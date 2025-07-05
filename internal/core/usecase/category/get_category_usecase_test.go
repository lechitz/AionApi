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

func TestGetCategoryByID_InvalidCategoryID(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByID(suite.Ctx, category).
		Return(domain.Category{}, errors.New(constants.FailedToGetCategoryByID))

	categoryDB, err := suite.CategoryService.GetCategoryByID(suite.Ctx, category)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, categoryDB)
	require.Equal(t, constants.FailedToGetCategoryByID, err.Error())
}

func TestGetCategoryByID_ErrorToGetCategoryByID(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByID(suite.Ctx, category).
		Return(domain.Category{}, errors.New(constants.FailedToGetCategoryByID))

	_, err := suite.CategoryService.GetCategoryByID(suite.Ctx, category)

	require.Error(t, err)
	require.Equal(t, constants.FailedToGetCategoryByID, err.Error())
}

func TestGetCategoryByID_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.CategoryRepository.EXPECT().
		GetCategoryByID(suite.Ctx, testdata.PerfectCategory).
		Return(testdata.PerfectCategory, nil)

	categoryDB, err := suite.CategoryService.GetCategoryByID(suite.Ctx, testdata.PerfectCategory)

	require.NoError(t, err)
	require.NotNil(t, categoryDB)
	require.Equal(t, testdata.PerfectCategory.ID, categoryDB.ID)
}

func TestGetCategoryByName_InvalidCategoryName(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByName(suite.Ctx, category).
		Return(domain.Category{}, errors.New(constants.CategoryNameIsRequired))

	categoryDB, err := suite.CategoryService.GetCategoryByName(suite.Ctx, category)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, categoryDB)
	require.Equal(t, constants.CategoryNameIsRequired, err.Error())
}

func TestGetCategoryByName_ErrorToGetCategoryByName(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByName(suite.Ctx, category).
		Return(domain.Category{}, errors.New(constants.FailedToGetCategoryByName))

	categoryDB, err := suite.CategoryService.GetCategoryByName(suite.Ctx, category)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, categoryDB)
	require.Equal(t, constants.FailedToGetCategoryByName, err.Error())
}

func TestGetCategoryByName_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByName(suite.Ctx, category).
		Return(testdata.PerfectCategory, nil)

	categoryDB, err := suite.CategoryService.GetCategoryByName(suite.Ctx, category)

	require.NoError(t, err)
	require.NotNil(t, categoryDB)
	require.Equal(t, testdata.PerfectCategory.ID, categoryDB.ID)
}

func TestGetAllCategories_ErrorToGetAllCategories(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := testdata.TestPerfectUser.ID

	suite.CategoryRepository.EXPECT().
		GetAllCategories(suite.Ctx, userID).
		Return(nil, errors.New(constants.FailedToGetAllCategories))

	categories, err := suite.CategoryService.GetAllCategories(suite.Ctx, userID)

	require.Error(t, err)
	require.Nil(t, categories)
	require.Equal(t, constants.FailedToGetAllCategories, err.Error())
}

func TestGetAllCategories_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := testdata.TestPerfectUser.ID

	suite.CategoryRepository.EXPECT().
		GetAllCategories(suite.Ctx, userID).
		Return([]domain.Category{testdata.PerfectCategory}, nil)

	categories, err := suite.CategoryService.GetAllCategories(suite.Ctx, userID)

	require.NoError(t, err)
	require.NotNil(t, categories)
	require.Len(t, categories, 1)
	require.Equal(t, testdata.PerfectCategory, categories[0])
}
