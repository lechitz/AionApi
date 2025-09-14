package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/category/core/domain"
	"github.com/lechitz/AionApi/internal/category/core/ports/input"
	"github.com/lechitz/AionApi/internal/category/core/usecase"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/lechitz/AionApi/tests/testdata"
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

	c := testdata.PerfectCategory
	cmd := makeUpdateCmdFromDomain(c)

	updateFields := map[string]interface{}{
		commonkeys.CategoryName:        c.Name,
		commonkeys.CategoryDescription: c.Description,
		commonkeys.CategoryColor:       c.Color,
		commonkeys.CategoryIcon:        c.Icon,
	}

	suite.CategoryRepository.EXPECT().
		UpdateCategory(gomock.Any(), c.ID, c.UserID, updateFields).
		Return(domain.Category{}, errors.New(usecase.FailedToUpdateCategory))

	categoryDB, err := suite.CategoryService.Update(suite.Ctx, cmd)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, categoryDB)
	require.Equal(t, usecase.FailedToUpdateCategory, err.Error())
}

func TestUpdateCategory_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	c := testdata.PerfectCategory
	cmd := makeUpdateCmdFromDomain(c)

	updateFields := map[string]interface{}{
		commonkeys.CategoryName:        c.Name,
		commonkeys.CategoryDescription: c.Description,
		commonkeys.CategoryColor:       c.Color,
		commonkeys.CategoryIcon:        c.Icon,
	}

	suite.CategoryRepository.EXPECT().
		UpdateCategory(gomock.Any(), c.ID, c.UserID, updateFields).
		Return(c, nil)

	categoryDB, err := suite.CategoryService.Update(suite.Ctx, cmd)

	require.NoError(t, err)
	require.Equal(t, c.ID, categoryDB.ID)
	require.Equal(t, c.UserID, categoryDB.UserID)
	require.Equal(t, c.Name, categoryDB.Name)
	require.Equal(t, c.Description, categoryDB.Description)
	require.Equal(t, c.Color, categoryDB.Color)
	require.Equal(t, c.Icon, categoryDB.Icon)
}
