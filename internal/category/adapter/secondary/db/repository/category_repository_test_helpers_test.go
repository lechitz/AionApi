package repository_test

import (
	"testing"
	"time"

	repository "github.com/lechitz/aion-api/internal/category/adapter/secondary/db/repository"
	"github.com/lechitz/aion-api/internal/category/core/domain"
	"github.com/lechitz/aion-api/tests/mocks"
	"go.uber.org/mock/gomock"
)

func newCategoryRepo(t *testing.T) (*repository.CategoryRepository, *mocks.MockDB) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	dbMock := mocks.NewMockDB(ctrl)
	loggerMock := mocks.NewMockContextLogger(ctrl)
	loggerMock.EXPECT().Errorw(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	return repository.New(dbMock, loggerMock), dbMock
}

func sampleCategory() domain.Category {
	now := time.Now().UTC()
	return domain.Category{
		ID:          7,
		UserID:      42,
		Name:        "health",
		Description: "desc",
		Color:       "#00AA00",
		Icon:        "leaf",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
