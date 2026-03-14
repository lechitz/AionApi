package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/tag/core/domain"
	"github.com/lechitz/AionApi/internal/tag/core/usecase"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAll_InvalidUserID(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	tags, err := suite.TagService.GetAll(suite.Ctx, 0)

	require.Error(t, err)
	require.Equal(t, usecase.UserIDIsRequired, err.Error())
	require.Empty(t, tags)
}

func TestGetAll_CacheHit(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(100)
	expectedTags := []domain.Tag{
		{ID: 1, Name: "Work", UserID: userID},
		{ID: 2, Name: "Personal", UserID: userID},
	}

	suite.TagCache.EXPECT().
		GetTagList(gomock.Any(), userID).
		Return(expectedTags, nil)

	tags, err := suite.TagService.GetAll(suite.Ctx, userID)

	require.NoError(t, err)
	require.Equal(t, expectedTags, tags)
}

func TestGetAll_RepositoryError(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(100)

	suite.TagCache.EXPECT().
		GetTagList(gomock.Any(), userID).
		Return(nil, errors.New("cache miss"))

	suite.TagRepository.EXPECT().
		GetAll(gomock.Any(), userID).
		Return(nil, errors.New(usecase.FailedToListTags))

	tags, err := suite.TagService.GetAll(suite.Ctx, userID)

	require.Error(t, err)
	require.Empty(t, tags)
	require.Equal(t, usecase.FailedToListTags, err.Error())
}

func TestGetAll_Success(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(100)
	expectedTags := []domain.Tag{
		{ID: 1, Name: "Work", UserID: userID},
		{ID: 2, Name: "Personal", UserID: userID},
		{ID: 3, Name: "Urgent", UserID: userID},
	}

	suite.TagCache.EXPECT().
		GetTagList(gomock.Any(), userID).
		Return(nil, errors.New("cache miss"))

	suite.TagRepository.EXPECT().
		GetAll(gomock.Any(), userID).
		Return(expectedTags, nil)

	suite.TagCache.EXPECT().
		SaveTagList(gomock.Any(), userID, expectedTags, gomock.Any()).
		Return(nil)

	tags, err := suite.TagService.GetAll(suite.Ctx, userID)

	require.NoError(t, err)
	require.Equal(t, expectedTags, tags)
}
