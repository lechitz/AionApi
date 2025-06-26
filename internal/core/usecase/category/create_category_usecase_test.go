package category_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain/entity"

	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/lechitz/AionApi/tests/testdata"
	"github.com/stretchr/testify/require"
)

func TestCreateCategory_ErrorToValidateCreateCategoryRequired_Name(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory
	category.Name = ""

	createdCategory, err := suite.CategoryService.CreateCategory(suite.Ctx, category)

	require.Error(t, err)
	require.Equal(t, entity.Category{}, createdCategory)
}

func TestCreateCategory_ErrorToValidateCreateCategoryRequired_DescriptionExceedLimit(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory
	category.Description = "Lorem ipsum lor sit amet, consectetur adipiscing elit. Donec nec justo eget felis facilisis. Aliquam porttitor mauris sit a orci. Aenean dignissim, orci non feugiat placerat, augue justo posuere justo."

	createdCategory, err := suite.CategoryService.CreateCategory(suite.Ctx, category)

	require.Error(t, err)
	require.Equal(t, entity.Category{}, createdCategory)
}

func TestCreateCategory_ErrorToGetCategoryByName(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByName(suite.Ctx, category).
		Return(category, nil)

	createdCategory, err := suite.CategoryService.CreateCategory(suite.Ctx, category)

	require.Error(t, err)
	require.Equal(t, entity.Category{}, createdCategory)
}

func TestCreateCategory_ErrorToCreateCategory(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByName(suite.Ctx, category).
		Return(entity.Category{}, nil)

	suite.CategoryRepository.EXPECT().
		CreateCategory(suite.Ctx, category).
		Return(entity.Category{}, errors.New(constants.CategoryAlreadyExists))

	createdCategory, err := suite.CategoryService.CreateCategory(suite.Ctx, category)

	require.Error(t, err)
	require.Equal(t, entity.Category{}, createdCategory)
}

func TestCreateCategory_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByName(suite.Ctx, category).
		Return(entity.Category{}, nil)

	suite.CategoryRepository.EXPECT().
		CreateCategory(suite.Ctx, category).
		Return(category, nil)

	createdCategory, err := suite.CategoryService.CreateCategory(suite.Ctx, category)

	require.NoError(t, err)
	require.Equal(t, category, createdCategory)
}
