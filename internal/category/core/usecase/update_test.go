package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/category/core/domain"
	"github.com/lechitz/AionApi/internal/category/core/ports/input"
	"github.com/lechitz/AionApi/internal/category/core/usecase"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// build an UpdateCategoryCommand from a domain.Category (used only in tests).
func makeUpdateCmdFromDomain(d domain.Category) input.UpdateCategoryCommand {
	var name, desc, color, icon *string
	if d.Name != "" {
		name = &d.Name
	}
	if d.Description != "" {
		desc = &d.Description
	}
	if d.Color != "" {
		color = &d.Color
	}
	if d.Icon != "" {
		icon = &d.Icon
	}
	return input.UpdateCategoryCommand{
		ID:          d.ID,
		UserID:      d.UserID,
		Name:        name,
		Description: desc,
		ColorHex:    color,
		Icon:        icon,
	}
}

func TestUpdateCategory_ErrorToUpdateCategory(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := domain.Category{
		ID:          1,
		UserID:      3,
		Name:        "Work",
		Description: "my work description",
		Color:       "blue",
		Icon:        "work/briefcase.svg",
	}
	cmd := makeUpdateCmdFromDomain(category)

	updateFields := map[string]interface{}{
		commonkeys.CategoryName:        category.Name,
		commonkeys.CategoryDescription: category.Description,
		commonkeys.CategoryColor:       category.Color,
		commonkeys.CategoryIcon:        category.Icon,
	}

	suite.CategoryRepository.EXPECT().
		GetByName(gomock.Any(), category.Name, category.UserID).
		Return(domain.Category{}, nil)

	suite.CategoryRepository.EXPECT().
		UpdateCategory(gomock.Any(), category.ID, category.UserID, updateFields).
		Return(domain.Category{}, errors.New(usecase.FailedToUpdateCategory))

	categoryDB, err := suite.CategoryService.Update(suite.Ctx, cmd)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, categoryDB)
	require.Equal(t, usecase.FailedToUpdateCategory, err.Error())
}

func TestUpdateCategory_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	category := domain.Category{
		ID:          1,
		UserID:      3,
		Name:        "Work",
		Description: "my work description",
		Color:       "blue",
		Icon:        "work/briefcase.svg",
	}
	cmd := makeUpdateCmdFromDomain(category)

	updateFields := map[string]interface{}{
		commonkeys.CategoryName:        category.Name,
		commonkeys.CategoryDescription: category.Description,
		commonkeys.CategoryColor:       category.Color,
		commonkeys.CategoryIcon:        category.Icon,
	}

	suite.CategoryRepository.EXPECT().
		GetByName(gomock.Any(), category.Name, category.UserID).
		Return(domain.Category{}, nil)

	suite.CategoryRepository.EXPECT().
		UpdateCategory(gomock.Any(), category.ID, category.UserID, updateFields).
		Return(category, nil)

	suite.CategoryCache.EXPECT().
		DeleteCategory(gomock.Any(), category.ID, category.UserID).
		Return(nil)

	suite.CategoryCache.EXPECT().
		DeleteCategoryByName(gomock.Any(), category.Name, category.UserID).
		Return(nil)

	suite.CategoryCache.EXPECT().
		DeleteCategoryList(gomock.Any(), category.UserID).
		Return(nil)

	categoryDB, err := suite.CategoryService.Update(suite.Ctx, cmd)

	require.NoError(t, err)
	require.Equal(t, category.ID, categoryDB.ID)
	require.Equal(t, category.UserID, categoryDB.UserID)
	require.Equal(t, category.Name, categoryDB.Name)
	require.Equal(t, category.Description, categoryDB.Description)
	require.Equal(t, category.Color, categoryDB.Color)
	require.Equal(t, category.Icon, categoryDB.Icon)
}

func TestUpdateCategory_NameAlreadyExists(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	targetID := uint64(1)
	differentID := uint64(2)
	userID := uint64(3)
	name := "Work"

	cmd := input.UpdateCategoryCommand{
		ID:     targetID,
		UserID: userID,
		Name:   &name,
	}

	suite.CategoryRepository.EXPECT().
		GetByName(gomock.Any(), name, userID).
		Return(domain.Category{ID: differentID, UserID: userID, Name: name}, nil)

	updated, err := suite.CategoryService.Update(suite.Ctx, cmd)
	require.Error(t, err)
	require.Equal(t, usecase.CategoryAlreadyExists, err.Error())
	require.Equal(t, domain.Category{}, updated)
}

func TestUpdateCategory_ErrorToValidateCategory_IconInvalid(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	invalidIcon := "work.png"
	cmd := input.UpdateCategoryCommand{
		ID:     1,
		UserID: 3,
		Icon:   &invalidIcon,
	}

	updated, err := suite.CategoryService.Update(suite.Ctx, cmd)
	require.Error(t, err)
	require.Equal(t, usecase.CategoryIconInvalid, err.Error())
	require.Equal(t, domain.Category{}, updated)
}
