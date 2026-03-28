package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/aion-api/internal/category/core/domain"
	"github.com/lechitz/aion-api/internal/category/core/usecase"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListAll_ErrorToGetAllCategories(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(3)

	suite.CategoryCache.EXPECT().
		GetCategoryList(gomock.Any(), userID).
		Return(nil, errors.New("cache miss"))

	suite.CategoryRepository.EXPECT().
		ListAll(gomock.Any(), userID).
		Return(nil, errors.New(usecase.FailedToGetAllCategories))

	categories, err := suite.CategoryService.ListAll(suite.Ctx, userID)

	require.Error(t, err)
	require.Nil(t, categories)
	require.Equal(t, usecase.FailedToGetAllCategories, err.Error())
}

func TestListAll_Success(t *testing.T) {
	suite := setup.CategoryServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(3)
	expected := []domain.Category{
		{
			ID:     1,
			UserID: 3,
			Name:   "Work",
		},
	}

	suite.CategoryCache.EXPECT().
		GetCategoryList(gomock.Any(), userID).
		Return(nil, errors.New("cache miss"))

	suite.CategoryRepository.EXPECT().
		ListAll(gomock.Any(), userID).
		Return(expected, nil)

	suite.CategoryCache.EXPECT().
		SaveCategoryList(gomock.Any(), userID, expected, gomock.Any()).
		Return(nil)

	categories, err := suite.CategoryService.ListAll(suite.Ctx, userID)

	require.NoError(t, err)
	require.Equal(t, expected, categories)
}
