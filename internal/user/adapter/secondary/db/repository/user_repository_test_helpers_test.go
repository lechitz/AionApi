package repository_test

import (
	"testing"
	"time"

	repository "github.com/lechitz/aion-api/internal/user/adapter/secondary/db/repository"
	"github.com/lechitz/aion-api/internal/user/core/domain"
	"github.com/lechitz/aion-api/tests/mocks"
	"go.uber.org/mock/gomock"
)

func newUserRepo(t *testing.T) (*repository.UserRepository, *mocks.MockDB, *mocks.MockRoleAssigner) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	dbMock := mocks.NewMockDB(ctrl)
	loggerMock := mocks.NewMockContextLogger(ctrl)
	assignerMock := mocks.NewMockRoleAssigner(ctrl)

	loggerMock.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().InfowCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().ErrorwCtx(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	return repository.New(dbMock, loggerMock, assignerMock), dbMock, assignerMock
}

func sampleUser() domain.User {
	now := time.Now().UTC()
	return domain.User{
		ID:        5,
		Name:      "John",
		Username:  "john",
		Email:     "john@example.com",
		Password:  "hash",
		CreatedAt: now,
		UpdatedAt: now,
	}
}
