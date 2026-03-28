package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/aion-api/internal/tag/core/domain"
	"github.com/lechitz/aion-api/internal/tag/core/usecase"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetByID_InvalidTagID(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	tag, err := suite.TagService.GetByID(suite.Ctx, 0, 1)

	require.Error(t, err)
	require.Equal(t, usecase.UserIDIsRequired, err.Error())
	require.Equal(t, domain.Tag{}, tag)
}

func TestGetByID_InvalidUserID(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	tag, err := suite.TagService.GetByID(suite.Ctx, 1, 0)

	require.Error(t, err)
	require.Equal(t, usecase.UserIDIsRequired, err.Error())
	require.Equal(t, domain.Tag{}, tag)
}

func TestGetByID_CacheHit(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	expectedTag := domain.Tag{
		ID:         1,
		Name:       "Work",
		CategoryID: 10,
		UserID:     100,
	}

	suite.TagCache.EXPECT().
		GetTag(gomock.Any(), expectedTag.ID, expectedTag.UserID).
		Return(expectedTag, nil)

	tag, err := suite.TagService.GetByID(suite.Ctx, expectedTag.ID, expectedTag.UserID)

	require.NoError(t, err)
	require.Equal(t, expectedTag, tag)
}

func TestGetByID_RepositoryError(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	tagID := uint64(999)
	userID := uint64(100)

	suite.TagCache.EXPECT().
		GetTag(gomock.Any(), tagID, userID).
		Return(domain.Tag{}, errors.New("cache miss"))

	suite.TagRepository.EXPECT().
		GetByID(gomock.Any(), tagID, userID).
		Return(domain.Tag{}, errors.New("tag not found"))

	tag, err := suite.TagService.GetByID(suite.Ctx, tagID, userID)

	require.Error(t, err)
	require.Equal(t, domain.Tag{}, tag)
}

func TestGetByID_Success(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	expectedTag := domain.Tag{
		ID:         1,
		Name:       "Work",
		CategoryID: 10,
		UserID:     100,
	}

	suite.TagCache.EXPECT().
		GetTag(gomock.Any(), expectedTag.ID, expectedTag.UserID).
		Return(domain.Tag{}, errors.New("cache miss"))

	suite.TagRepository.EXPECT().
		GetByID(gomock.Any(), expectedTag.ID, expectedTag.UserID).
		Return(expectedTag, nil)

	suite.TagCache.EXPECT().
		SaveTag(gomock.Any(), expectedTag, gomock.Any()).
		Return(nil)

	tag, err := suite.TagService.GetByID(suite.Ctx, expectedTag.ID, expectedTag.UserID)

	require.NoError(t, err)
	require.Equal(t, expectedTag, tag)
}
