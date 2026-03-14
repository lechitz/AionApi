package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/tag/core/domain"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetByCategoryID_CacheHit(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	categoryID := uint64(10)
	userID := uint64(100)
	expectedTags := []domain.Tag{
		{ID: 1, Name: "Work", CategoryID: categoryID, UserID: userID},
		{ID: 2, Name: "Important", CategoryID: categoryID, UserID: userID},
	}

	suite.TagCache.EXPECT().
		GetTagsByCategory(gomock.Any(), categoryID, userID).
		Return(expectedTags, nil)

	tags, err := suite.TagService.GetByCategoryID(suite.Ctx, categoryID, userID)

	require.NoError(t, err)
	require.Equal(t, expectedTags, tags)
}

func TestGetByCategoryID_RepositoryError(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	categoryID := uint64(999)
	userID := uint64(100)

	suite.TagCache.EXPECT().
		GetTagsByCategory(gomock.Any(), categoryID, userID).
		Return(nil, errors.New("cache miss"))

	suite.TagRepository.EXPECT().
		GetByCategoryID(gomock.Any(), categoryID, userID).
		Return(nil, errors.New("failed to get tags"))

	tags, err := suite.TagService.GetByCategoryID(suite.Ctx, categoryID, userID)

	require.Error(t, err)
	require.Empty(t, tags)
}

func TestGetByCategoryID_Success(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	categoryID := uint64(10)
	userID := uint64(100)
	expectedTags := []domain.Tag{
		{ID: 1, Name: "Work", CategoryID: categoryID, UserID: userID},
		{ID: 2, Name: "Important", CategoryID: categoryID, UserID: userID},
		{ID: 3, Name: "Urgent", CategoryID: categoryID, UserID: userID},
	}

	suite.TagCache.EXPECT().
		GetTagsByCategory(gomock.Any(), categoryID, userID).
		Return(nil, errors.New("cache miss"))

	suite.TagRepository.EXPECT().
		GetByCategoryID(gomock.Any(), categoryID, userID).
		Return(expectedTags, nil)

	suite.TagCache.EXPECT().
		SaveTagsByCategory(gomock.Any(), categoryID, userID, expectedTags, gomock.Any()).
		Return(nil)

	tags, err := suite.TagService.GetByCategoryID(suite.Ctx, categoryID, userID)

	require.NoError(t, err)
	require.Equal(t, expectedTags, tags)
}

func TestGetByCategoryID_EmptyResult(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	categoryID := uint64(999)
	userID := uint64(100)

	suite.TagCache.EXPECT().
		GetTagsByCategory(gomock.Any(), categoryID, userID).
		Return(nil, errors.New("cache miss"))

	suite.TagRepository.EXPECT().
		GetByCategoryID(gomock.Any(), categoryID, userID).
		Return([]domain.Tag{}, nil)

	suite.TagCache.EXPECT().
		SaveTagsByCategory(gomock.Any(), categoryID, userID, []domain.Tag{}, gomock.Any()).
		Return(nil)

	tags, err := suite.TagService.GetByCategoryID(suite.Ctx, categoryID, userID)

	require.NoError(t, err)
	require.Empty(t, tags)
}
