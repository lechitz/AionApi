package usecase_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/lechitz/AionApi/internal/category/core/domain"
	"github.com/lechitz/AionApi/internal/category/core/ports/input"
	"github.com/lechitz/AionApi/internal/category/core/usecase"
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

	cmd := makeCreateCmdFromDomain(category)

	createdCategory, err := suite.CategoryService.Create(suite.Ctx, cmd)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, createdCategory)
}

func TestCreateCategory_ErrorToValidateCreateCategoryRequired_DescriptionExceedLimit(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory
	category.Description = strings.Repeat("x", 201)

	cmd := makeCreateCmdFromDomain(category)

	createdCategory, err := suite.CategoryService.Create(suite.Ctx, cmd)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, createdCategory)
}

// Test for ColorHex length exceeding the maximum allowed.
func TestCreateCategory_ErrorToValidateCreateCategoryRequired_ColorExceedLimit(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory
	category.Color = "#12345678" // length > 7

	cmd := makeCreateCmdFromDomain(category)

	createdCategory, err := suite.CategoryService.Create(suite.Ctx, cmd)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, createdCategory)
}

// Test for Icon length exceeding the maximum allowed.
func TestCreateCategory_ErrorToValidateCreateCategoryRequired_IconExceedLimit(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory
	category.Icon = strings.Repeat("i", 51) // length > 50

	cmd := makeCreateCmdFromDomain(category)

	createdCategory, err := suite.CategoryService.Create(suite.Ctx, cmd)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, createdCategory)
}

func TestCreateCategory_ErrorToGetCategoryByName(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory
	cmd := makeCreateCmdFromDomain(category)

	suite.CategoryRepository.EXPECT().
		GetByName(gomock.Any(), category.Name, category.UserID).
		Return(category, nil)

	createdCategory, err := suite.CategoryService.Create(suite.Ctx, cmd)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, createdCategory)
}

func TestCreateCategory_ErrorToCreateCategory(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory
	cmd := makeCreateCmdFromDomain(category)

	suite.CategoryRepository.EXPECT().
		GetByName(gomock.Any(), category.Name, category.UserID).
		Return(domain.Category{}, nil)

	suite.CategoryRepository.EXPECT().
		Create(gomock.Any(), domain.Category{
			UserID:      category.UserID,
			Name:        category.Name,
			Description: category.Description,
			Color:       category.Color,
			Icon:        category.Icon,
		}).
		Return(domain.Category{}, errors.New(usecase.FailedToCreateCategory))

	createdCategory, err := suite.CategoryService.Create(suite.Ctx, cmd)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, createdCategory)
}

// Test that nil pointers result in empty strings.
func TestCreateCategory_PtrOrEmpty_NilPointersBecomeEmptyStrings(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := domain.Category{
		UserID: 3,
		Name:   "Work",
	}

	cmd := makeCreateCmdFromDomain(category)

	suite.CategoryRepository.EXPECT().
		GetByName(gomock.Any(), category.Name, category.UserID).
		Return(domain.Category{}, nil)

	suite.CategoryRepository.EXPECT().
		Create(gomock.Any(), domain.Category{
			UserID:      category.UserID,
			Name:        category.Name,
			Description: "",
			Color:       "",
			Icon:        "",
		}).
		Return(domain.Category{
			UserID:      category.UserID,
			Name:        category.Name,
			Description: "",
			Color:       "",
			Icon:        "",
		}, nil)

	created, err := suite.CategoryService.Create(suite.Ctx, cmd)
	require.NoError(t, err)
	require.Empty(t, created.Description)
	require.Empty(t, created.Color)
	require.Empty(t, created.Icon)
}

func TestCreateCategory_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := testdata.PerfectCategory
	cmd := makeCreateCmdFromDomain(category)

	suite.CategoryRepository.EXPECT().
		GetByName(gomock.Any(), category.Name, category.UserID).
		Return(domain.Category{}, nil)

	suite.CategoryRepository.EXPECT().
		Create(gomock.Any(), domain.Category{
			UserID:      category.UserID,
			Name:        category.Name,
			Description: category.Description,
			Color:       category.Color,
			Icon:        category.Icon,
		}).
		Return(category, nil)

	createdCategory, err := suite.CategoryService.Create(suite.Ctx, cmd)

	require.NoError(t, err)
	require.Equal(t, category, createdCategory)
}

// build a CreateCategoryCommand from a domain.Category (used only in tests).
func makeCreateCmdFromDomain(d domain.Category) input.CreateCategoryCommand {
	var desc, color, icon *string
	if d.Description != "" {
		desc = &d.Description
	}
	if d.Color != "" {
		color = &d.Color
	}
	if d.Icon != "" {
		icon = &d.Icon
	}
	return input.CreateCategoryCommand{
		Name:        d.Name,
		Description: desc,
		ColorHex:    color,
		Icon:        icon,
		UserID:      d.UserID,
	}
}
