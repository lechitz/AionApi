package usecase_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/lechitz/AionApi/internal/tag/core/domain"
	"github.com/lechitz/AionApi/internal/tag/core/ports/input"
	"github.com/lechitz/AionApi/internal/tag/core/usecase"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// ---------- Helpers ----------

func perfectTag() domain.Tag {
	return domain.Tag{
		UserID:      1,
		CategoryID:  2,
		Name:        "Read",
		Description: "Daily reading practice",
	}
}

func makeCreateTagCmdFromDomain(d domain.Tag) input.CreateTagCommand {
	var desc *string
	if d.Description != "" {
		desc = &d.Description
	}
	return input.CreateTagCommand{
		Name:        d.Name,
		UserID:      d.UserID,
		CategoryID:  d.CategoryID,
		Description: desc,
	}
}

// ---------- Tests ----------

func TestCreateTag_ErrorToValidateCreateTagRequired_Name(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	tag := perfectTag()
	tag.Name = ""

	cmd := makeCreateTagCmdFromDomain(tag)

	created, err := suite.TagService.Create(suite.Ctx, cmd)

	require.Error(t, err)
	require.Equal(t, domain.Tag{}, created)
}

func TestCreateTag_ErrorToValidateCreateTagRequired_DescriptionExceedLimit(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	tag := perfectTag()
	tag.Description = strings.Repeat("x", 201) // > 200

	cmd := makeCreateTagCmdFromDomain(tag)

	created, err := suite.TagService.Create(suite.Ctx, cmd)

	require.Error(t, err)
	require.Equal(t, domain.Tag{}, created)
}

func TestCreateTag_ErrorToGetTagByName(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	tag := perfectTag()
	cmd := makeCreateTagCmdFromDomain(tag)

	// Simula "já existe" (GetByName retorna um registro válido).
	suite.TagRepository.EXPECT().
		GetByName(gomock.Any(), tag.Name, tag.UserID).
		Return(tag, nil)

	created, err := suite.TagService.Create(suite.Ctx, cmd)

	require.Error(t, err)
	require.Equal(t, domain.Tag{}, created)
}

func TestCreateTag_ErrorToCreateTag(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	tag := perfectTag()
	cmd := makeCreateTagCmdFromDomain(tag)

	suite.TagRepository.EXPECT().
		GetByName(gomock.Any(), tag.Name, tag.UserID).
		Return(domain.Tag{}, nil)

	suite.TagRepository.EXPECT().
		Create(gomock.Any(), domain.Tag{
			UserID:      tag.UserID,
			CategoryID:  tag.CategoryID,
			Name:        tag.Name,
			Description: tag.Description,
		}).
		Return(domain.Tag{}, errors.New(usecase.FailedToCreateTag))

	created, err := suite.TagService.Create(suite.Ctx, cmd)

	require.Error(t, err)
	require.Equal(t, domain.Tag{}, created)
}

func TestCreateTag_PtrOrEmpty_NilPointersBecomeEmptyStrings(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	// Description vazia => ponteiro nil no cmd
	tag := domain.Tag{
		UserID:     10,
		CategoryID: 20,
		Name:       "Focus",
		// Description == "" (nil no cmd)
	}

	cmd := makeCreateTagCmdFromDomain(tag)

	suite.TagRepository.EXPECT().
		GetByName(gomock.Any(), tag.Name, tag.UserID).
		Return(domain.Tag{}, nil)

	// O usecase monta domain.Tag com Description "" (ptrOrEmpty)
	expectedCreate := domain.Tag{
		UserID:      tag.UserID,
		CategoryID:  tag.CategoryID,
		Name:        tag.Name,
		Description: "",
	}

	suite.TagRepository.EXPECT().
		Create(gomock.Any(), expectedCreate).
		Return(expectedCreate, nil)

	created, err := suite.TagService.Create(suite.Ctx, cmd)
	require.NoError(t, err)
	require.Equal(t, "", created.Description)
}

func TestCreateTag_Success(t *testing.T) {
	suite := setup.TagServiceTest(t)
	defer suite.Ctrl.Finish()

	tag := perfectTag()
	cmd := makeCreateTagCmdFromDomain(tag)

	suite.TagRepository.EXPECT().
		GetByName(gomock.Any(), tag.Name, tag.UserID).
		Return(domain.Tag{}, nil)

	suite.TagRepository.EXPECT().
		Create(gomock.Any(), domain.Tag{
			UserID:      tag.UserID,
			CategoryID:  tag.CategoryID,
			Name:        tag.Name,
			Description: tag.Description,
		}).
		Return(tag, nil)

	created, err := suite.TagService.Create(suite.Ctx, cmd)

	require.NoError(t, err)
	require.Equal(t, tag, created)
}
