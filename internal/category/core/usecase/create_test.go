package usecase_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/lechitz/AionApi/internal/core/category/domain"
	"github.com/lechitz/AionApi/internal/core/category/usecase"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/lechitz/AionApi/tests/testdata"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateCategory_ErrorToValidateCreateCategoryRequired_Name(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory
	category.Name = ""

	createdCategory, err := suite.CategoryService.Create(suite.Ctx, category)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, createdCategory)
}

func TestCreateCategory_ErrorToValidateCreateCategoryRequired_DescriptionExceedLimit(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory
	category.Description = strings.Repeat("x", 201)

	createdCategory, err := suite.CategoryService.Create(suite.Ctx, category)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, createdCategory)
}

func TestCreateCategory_ErrorToGetCategoryByName(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetByName(gomock.Any(), category.Name).
		Return(category, nil)

	createdCategory, err := suite.CategoryService.Create(suite.Ctx, category)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, createdCategory)
}

func TestCreateCategory_ErrorToCreateCategory(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetByName(gomock.Any(), category.Name).
		Return(domain.Category{}, nil)

	suite.CategoryRepository.EXPECT().
		Create(gomock.Any(), category).
		Return(domain.Category{}, errors.New(usecase.FailedToCreateCategory))

	createdCategory, err := suite.CategoryService.Create(suite.Ctx, category)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, createdCategory)
}

func TestCreateCategory_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory

	suite.CategoryRepository.EXPECT().
		GetByName(gomock.Any(), category.Name).
		Return(domain.Category{}, nil)

	suite.CategoryRepository.EXPECT().
		Create(gomock.Any(), category).
		Return(category, nil)

	createdCategory, err := suite.CategoryService.Create(suite.Ctx, category)

	require.NoError(t, err)
	require.Equal(t, category, createdCategory)
}
