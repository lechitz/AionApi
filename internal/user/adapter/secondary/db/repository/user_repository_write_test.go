package repository_test

import (
	"errors"
	"testing"

	"github.com/lechitz/aion-api/internal/platform/ports/output/db"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	usermodel "github.com/lechitz/aion-api/internal/user/adapter/secondary/db/model"
	repository "github.com/lechitz/aion-api/internal/user/adapter/secondary/db/repository"
	useroutput "github.com/lechitz/aion-api/internal/user/core/ports/output"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestCheckUniqueness(t *testing.T) {
	repo, dbMock, _ := newUserRepo(t)

	t.Run("username and email taken", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Select(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			row, ok := dest.(*usermodel.UserDB)
			require.True(t, ok)
			row.ID = 101
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Select(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			row, ok := dest.(*usermodel.UserDB)
			require.True(t, ok)
			row.ID = 202
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		res, err := repo.CheckUniqueness(t.Context(), "john", "john@example.com")
		require.NoError(t, err)
		require.Equal(t, useroutput.UserUniqueness{UsernameTaken: true, EmailTaken: true, UsernameOwnerID: ptr(101), EmailOwnerID: ptr(202)}, res)
	})

	t.Run("db error on username check", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Select(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Error().Return(errors.New("lookup fail"))

		_, err := repo.CheckUniqueness(t.Context(), "john", "")
		require.Error(t, err)
	})
}

func TestCreateUserRepository(t *testing.T) {
	repo, dbMock, assigner := newUserRepo(t)

	t.Run("unique violation returns validation error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Transaction(gomock.Any()).Return(errors.New("pq: duplicate key users_email_key"))

		_, err := repo.Create(t.Context(), sampleUser())
		require.Error(t, err)
		var vErr *sharederrors.ValidationError
		require.ErrorAs(t, err, &vErr)
	})

	t.Run("unique username violation returns validation error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Transaction(gomock.Any()).Return(errors.New("pq: duplicate key users_username_key"))

		_, err := repo.Create(t.Context(), sampleUser())
		require.Error(t, err)
		var vErr *sharederrors.ValidationError
		require.ErrorAs(t, err, &vErr)
	})

	t.Run("generic transaction error", func(t *testing.T) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Transaction(gomock.Any()).Return(errors.New("tx fail"))

		_, err := repo.Create(t.Context(), sampleUser())
		require.EqualError(t, err, "tx fail")
	})

	t.Run("success", func(t *testing.T) {
		assigner.EXPECT().AssignDefaultRole(gomock.Any(), uint64(77)).Return(nil)

		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Transaction(gomock.Any()).DoAndReturn(func(fc func(db.DB) error) error {
			tx := dbMock
			tx.EXPECT().Create(gomock.Any()).DoAndReturn(func(v any) db.DB {
				row, ok := v.(*usermodel.UserDB)
				require.True(t, ok)
				row.ID = 77
				return tx
			})
			tx.EXPECT().Error().Return(nil)
			return fc(tx)
		})

		// GetByID after create
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Model(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().First(gomock.Any()).DoAndReturn(func(dest any, _ ...any) db.DB {
			row, ok := dest.(*usermodel.UserDB)
			require.True(t, ok)
			row.ID = 77
			row.Username = "john"
			return dbMock
		})
		dbMock.EXPECT().Error().Return(nil)

		got, err := repo.Create(t.Context(), sampleUser())
		require.NoError(t, err)
		require.Equal(t, uint64(77), got.ID)
	})
}

func TestGetByUsernameGenericError(t *testing.T) {
	repo, dbMock, _ := newUserRepo(t)
	dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Select(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Where(gomock.Any(), gomock.Any()).Return(dbMock)
	dbMock.EXPECT().First(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Error().Return(errors.New("db fail"))

	_, err := repo.GetByUsername(t.Context(), "john")
	require.Error(t, err)
}

func ptr(v uint64) *uint64 { return &v }

var (
	_ = gorm.ErrRecordNotFound
	_ = repository.TracerName
)
