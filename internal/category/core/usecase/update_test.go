package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/category/domain"
	"github.com/lechitz/AionApi/internal/core/category/usecase"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.uber.org/mock/gomock"

	"github.com/lechitz/AionApi/tests/setup"
	"github.com/lechitz/AionApi/tests/testdata"
	"github.com/stretchr/testify/require"
)

func TestUpdateCategory_ErrorToUpdateCategory(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	c := testdata.PerfectCategory
	updateFields := map[string]interface{}{
		commonkeys.CategoryName:        c.Name,
		commonkeys.CategoryDescription: c.Description,
		commonkeys.CategoryColor:       c.Color,
		commonkeys.CategoryIcon:        c.Icon,
	}

	suite.CategoryRepository.EXPECT().
		UpdateCategory(gomock.Any(), c.ID, c.UserID, updateFields).
		Return(domain.Category{}, errors.New(usecase.FailedToUpdateCategory))

	categoryDB, err := suite.CategoryService.Update(suite.Ctx, c)

	require.Error(t, err)
	require.Equal(t, domain.Category{}, categoryDB)
	require.Equal(t, usecase.FailedToUpdateCategory, err.Error())
}

func TestUpdateCategory_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	c := testdata.PerfectCategory
	updateFields := map[string]interface{}{
		commonkeys.CategoryName:        c.Name,
		commonkeys.CategoryDescription: c.Description,
		commonkeys.CategoryColor:       c.Color,
		commonkeys.CategoryIcon:        c.Icon,
	}

	suite.CategoryRepository.EXPECT().
		UpdateCategory(gomock.Any(), c.ID, c.UserID, updateFields).
		Return(c, nil)

	categoryDB, err := suite.CategoryService.Update(suite.Ctx, c)

	require.NoError(t, err)
	require.Equal(t, c.ID, categoryDB.ID)
	require.Equal(t, c.UserID, categoryDB.UserID)
	require.Equal(t, c.Name, categoryDB.Name)
	require.Equal(t, c.Description, categoryDB.Description)
	require.Equal(t, c.Color, categoryDB.Color)
	require.Equal(t, c.Icon, categoryDB.Icon)
}
