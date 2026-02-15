package repository_test

import (
	"testing"
	"time"

	repository "github.com/lechitz/AionApi/internal/tag/adapter/secondary/db/repository"
	"github.com/lechitz/AionApi/internal/tag/core/domain"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/mock/gomock"
)

func newTagRepo(t *testing.T) (*repository.TagRepository, *mocks.MockDB) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	dbMock := mocks.NewMockDB(ctrl)
	loggerMock := mocks.NewMockContextLogger(ctrl)
	loggerMock.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	return repository.New(dbMock, loggerMock), dbMock
}

func sampleTag() domain.Tag {
	now := time.Now().UTC()
	return domain.Tag{
		ID:          10,
		UserID:      20,
		CategoryID:  30,
		Name:        "focus",
		Description: "tag desc",
		Icon:        "*",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
