package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/aion-api/internal/tag/core/usecase"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSoftDelete_InvalidTagID(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	err := suite.TagService.SoftDelete(suite.Ctx, 0, 100)

	require.Error(t, err)
	require.Equal(t, usecase.FailedToSoftDeleteTag, err.Error())
}

func TestSoftDelete_InvalidUserID(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	err := suite.TagService.SoftDelete(suite.Ctx, 1, 0)

	require.Error(t, err)
	require.Equal(t, usecase.FailedToSoftDeleteTag, err.Error())
}

func TestSoftDelete_RepositoryError(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	tagID := uint64(999)
	userID := uint64(100)

	suite.TagRepository.EXPECT().
		SoftDelete(gomock.Any(), tagID, userID).
		Return(errors.New("tag not found"))

	err := suite.TagService.SoftDelete(suite.Ctx, tagID, userID)

	require.Error(t, err)
	require.Equal(t, usecase.FailedToSoftDeleteTag, err.Error())
}

func TestSoftDelete_Success(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	tagID := uint64(1)
	userID := uint64(100)

	suite.TagRepository.EXPECT().
		SoftDelete(gomock.Any(), tagID, userID).
		Return(nil)

	suite.TagCache.EXPECT().
		DeleteTag(gomock.Any(), tagID, userID).
		Return(nil)

	suite.TagCache.EXPECT().
		DeleteTagList(gomock.Any(), userID).
		Return(nil)

	err := suite.TagService.SoftDelete(suite.Ctx, tagID, userID)

	require.NoError(t, err)
}

func TestSoftDelete_Success_CacheInvalidationFails(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	tagID := uint64(1)
	userID := uint64(100)

	suite.TagRepository.EXPECT().
		SoftDelete(gomock.Any(), tagID, userID).
		Return(nil)

	suite.TagCache.EXPECT().
		DeleteTag(gomock.Any(), tagID, userID).
		Return(errors.New("cache error"))

	suite.TagCache.EXPECT().
		DeleteTagList(gomock.Any(), userID).
		Return(errors.New("cache error"))

	// Should still succeed even if cache invalidation fails
	err := suite.TagService.SoftDelete(suite.Ctx, tagID, userID)

	require.NoError(t, err)
}
