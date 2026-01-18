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

func TestGetByName_Error_TagNameRequired(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	got, err := suite.TagService.GetByName(suite.Ctx, "", 1)

	require.Error(t, err)
	require.Equal(t, domain.Tag{}, got)
}

func TestGetByName_Error_RepositoryFailure(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	const userID uint64 = 1
	const tagName = "Read"

	repoErr := errors.New(usecase.FailedToGetTagByName)

	suite.TagCache.EXPECT().
		GetTagByName(gomock.Any(), tagName, userID).
		Return(domain.Tag{}, errors.New("cache miss"))

	suite.TagRepository.EXPECT().
		GetByName(gomock.Any(), tagName, userID).
		Return(domain.Tag{}, repoErr)

	got, err := suite.TagService.GetByName(suite.Ctx, tagName, userID)

	require.Error(t, err)
	require.Equal(t, domain.Tag{}, got)
}

func TestGetByName_Success_Found(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	const userID uint64 = 1
	const tagName = "Read"

	want := domain.Tag{
		ID:          42,
		UserID:      userID,
		CategoryID:  2,
		Name:        tagName,
		Description: "Daily reading practice",
	}

	suite.TagCache.EXPECT().
		GetTagByName(gomock.Any(), tagName, userID).
		Return(domain.Tag{}, errors.New("cache miss"))

	suite.TagRepository.EXPECT().
		GetByName(gomock.Any(), tagName, userID).
		Return(want, nil)

	suite.TagCache.EXPECT().
		SaveTagByName(gomock.Any(), want, gomock.Any()).
		Return(nil)

	suite.TagCache.EXPECT().
		SaveTag(gomock.Any(), want, gomock.Any()).
		Return(nil)

	got, err := suite.TagService.GetByName(suite.Ctx, tagName, userID)

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestGetByName_Success_NotFoundReturnsZeroValue(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	const userID uint64 = 1
	const tagName = "UnknownTag"

	suite.TagCache.EXPECT().
		GetTagByName(gomock.Any(), tagName, userID).
		Return(domain.Tag{}, errors.New("cache miss"))

	suite.TagRepository.EXPECT().
		GetByName(gomock.Any(), tagName, userID).
		Return(domain.Tag{}, nil)

	suite.TagCache.EXPECT().
		SaveTagByName(gomock.Any(), domain.Tag{}, gomock.Any()).
		Return(nil)

	got, err := suite.TagService.GetByName(suite.Ctx, tagName, userID)

	require.NoError(t, err)
	require.Equal(t, domain.Tag{}, got)
}
