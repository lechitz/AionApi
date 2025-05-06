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

func TestCreateCategory_ErrorToValidateCreateCategoryRequired_Name(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory
	category.Name = ""

	createdCategory, err := suite.CategoryService.CreateCategory(suite.Ctx, category)

	assert.Error(t, err)
	assert.Equal(t, domain.Category{}, createdCategory)
}

func TestCreateCategory_ErrorToValidateCreateCategoryRequired_DescriptionExceedLimit(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory
	category.Description = "Lorem ipsum lor sit amet, consectetur adipiscing elit. Donec nec justo eget felis facilisis. Aliquam porttitor mauris sit a orci. Aenean dignissim, orci non feugiat placerat, augue justo posuere justo."

	createdCategory, err := suite.CategoryService.CreateCategory(suite.Ctx, category)

	assert.Error(t, err)
	assert.Equal(t, domain.Category{}, createdCategory)
}

func TestCreateCategory_ErrorToGetCategoryByName(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByName(suite.Ctx, category).
		Return(category, nil)

	createdCategory, err := suite.CategoryService.CreateCategory(suite.Ctx, category)

	assert.Error(t, err)
	assert.Equal(t, domain.Category{}, createdCategory)
}

func TestCreateCategory_ErrorToCreateCategory(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByName(suite.Ctx, category).
		Return(domain.Category{}, nil)

	suite.CategoryRepository.EXPECT().
		CreateCategory(suite.Ctx, category).
		Return(domain.Category{}, errors.New(constants.CategoryAlreadyExists))

	createdCategory, err := suite.CategoryService.CreateCategory(suite.Ctx, category)

	assert.Error(t, err)
	assert.Equal(t, domain.Category{}, createdCategory)
}

func TestCreateCategory_Success(t *testing.T) {
	suite := setup.SetupCategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetCategoryByName(suite.Ctx, category).
		Return(domain.Category{}, nil)

	suite.CategoryRepository.EXPECT().
		CreateCategory(suite.Ctx, category).
		Return(category, nil)

	createdCategory, err := suite.CategoryService.CreateCategory(suite.Ctx, category)

	assert.NoError(t, err)
	assert.Equal(t, category, createdCategory)
}
