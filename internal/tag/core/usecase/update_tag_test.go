package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/tag/core/domain"
	"github.com/lechitz/AionApi/internal/tag/core/ports/input"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdate_Success(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	tagID := uint64(1)
	userID := uint64(100)
	newName := "Updated Work"
	newDesc := "Updated description"
	newCatID := uint64(20)

	cmd := input.UpdateTagCommand{
		ID:          tagID,
		UserID:      userID,
		Name:        &newName,
		Description: &newDesc,
		CategoryID:  &newCatID,
	}

	updateFields := map[string]interface{}{
		commonkeys.TagName:        newName,
		commonkeys.TagDescription: newDesc,
		commonkeys.CategoryID:     newCatID,
	}

	expectedTag := domain.Tag{
		ID:          tagID,
		Name:        newName,
		Description: newDesc,
		CategoryID:  newCatID,
		UserID:      userID,
	}

	suite.TagRepository.EXPECT().
		UpdateTag(gomock.Any(), tagID, userID, updateFields).
		Return(expectedTag, nil)

	suite.TagCache.EXPECT().
		DeleteTag(gomock.Any(), tagID, userID).
		Return(nil)

	suite.TagCache.EXPECT().
		DeleteTagByName(gomock.Any(), newName, userID).
		Return(nil)

	suite.TagCache.EXPECT().
		DeleteTagList(gomock.Any(), userID).
		Return(nil)

	suite.TagCache.EXPECT().
		DeleteTagsByCategory(gomock.Any(), newCatID, userID).
		Return(nil)

	tag, err := suite.TagService.Update(suite.Ctx, cmd)

	require.NoError(t, err)
	require.Equal(t, expectedTag, tag)
}

func TestUpdate_PartialUpdate_OnlyName(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	tagID := uint64(1)
	userID := uint64(100)
	newName := "Updated Name Only"

	cmd := input.UpdateTagCommand{
		ID:          tagID,
		UserID:      userID,
		Name:        &newName,
		Description: nil,
		CategoryID:  nil,
	}

	expectedTag := domain.Tag{
		ID:     tagID,
		Name:   newName,
		UserID: userID,
	}

	suite.TagRepository.EXPECT().
		UpdateTag(gomock.Any(), tagID, userID, gomock.Any()).
		DoAndReturn(func(_ any, _ uint64, _ uint64, fields map[string]interface{}) (domain.Tag, error) {
			require.Contains(t, fields, commonkeys.TagName)
			require.NotContains(t, fields, commonkeys.TagDescription)
			require.NotContains(t, fields, commonkeys.CategoryID)
			return expectedTag, nil
		})

	suite.TagCache.EXPECT().
		DeleteTag(gomock.Any(), tagID, userID).
		Return(nil)

	suite.TagCache.EXPECT().
		DeleteTagByName(gomock.Any(), newName, userID).
		Return(nil)

	suite.TagCache.EXPECT().
		DeleteTagList(gomock.Any(), userID).
		Return(nil)

	suite.TagCache.EXPECT().
		DeleteTagsByCategory(gomock.Any(), gomock.Any(), userID).
		Return(nil)

	tag, err := suite.TagService.Update(suite.Ctx, cmd)

	require.NoError(t, err)
	require.Equal(t, expectedTag, tag)
}

func TestUpdate_RepositoryError(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	tagID := uint64(999)
	userID := uint64(100)
	newName := "Non-existent"

	cmd := input.UpdateTagCommand{
		ID:     tagID,
		UserID: userID,
		Name:   &newName,
	}

	suite.TagRepository.EXPECT().
		UpdateTag(gomock.Any(), tagID, userID, gomock.Any()).
		Return(domain.Tag{}, errors.New("tag not found"))

	tag, err := suite.TagService.Update(suite.Ctx, cmd)

	require.Error(t, err)
	require.Equal(t, domain.Tag{}, tag)
}

func TestUpdate_Success_CacheInvalidationFails(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	tagID := uint64(1)
	userID := uint64(100)
	newName := "Updated"

	cmd := input.UpdateTagCommand{
		ID:     tagID,
		UserID: userID,
		Name:   &newName,
	}

	expectedTag := domain.Tag{
		ID:     tagID,
		Name:   newName,
		UserID: userID,
	}

	suite.TagRepository.EXPECT().
		UpdateTag(gomock.Any(), tagID, userID, gomock.Any()).
		Return(expectedTag, nil)

	// All cache invalidations fail - should still succeed
	suite.TagCache.EXPECT().
		DeleteTag(gomock.Any(), tagID, userID).
		Return(errors.New("cache error"))

	suite.TagCache.EXPECT().
		DeleteTagByName(gomock.Any(), newName, userID).
		Return(errors.New("cache error"))

	suite.TagCache.EXPECT().
		DeleteTagList(gomock.Any(), userID).
		Return(errors.New("cache error"))

	suite.TagCache.EXPECT().
		DeleteTagsByCategory(gomock.Any(), gomock.Any(), userID).
		Return(errors.New("cache error"))

	tag, err := suite.TagService.Update(suite.Ctx, cmd)

	require.NoError(t, err)
	require.Equal(t, expectedTag, tag)
}
