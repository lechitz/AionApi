package repository_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUserStatsWithMostUsedEntities(t *testing.T) {
	repo, dbMock, _ := newUserRepo(t)

	expectCountQuery := func(value int) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
		dbMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
			ptr, ok := dest.(*int)
			require.True(t, ok)
			*ptr = value
			return dbMock
		})
	}

	expectCountQuery(12)
	expectCountQuery(3)
	expectCountQuery(7)
	expectCountQuery(5)
	expectCountQuery(9)

	dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
		v := reflect.ValueOf(dest).Elem()
		v.FieldByName("CategoryID").SetUint(30)
		v.FieldByName("Name").SetString("work")
		v.FieldByName("Count").SetInt(4)
		return dbMock
	})
	dbMock.EXPECT().Error().Return(nil)

	dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest any) db.DB {
		v := reflect.ValueOf(dest).Elem()
		v.FieldByName("TagID").SetUint(88)
		v.FieldByName("Name").SetString("focus")
		v.FieldByName("Count").SetInt(11)
		return dbMock
	})
	dbMock.EXPECT().Error().Return(nil)

	got, err := repo.GetUserStats(t.Context(), 42)
	require.NoError(t, err)
	require.Equal(t, 12, got.TotalRecords)
	require.Equal(t, 3, got.TotalCategories)
	require.Equal(t, 7, got.TotalTags)
	require.Equal(t, 5, got.RecordsThisWeek)
	require.Equal(t, 9, got.RecordsThisMonth)
	require.NotNil(t, got.MostUsedCategory)
	require.Equal(t, uint64(30), got.MostUsedCategory.ID)
	require.NotNil(t, got.MostUsedTag)
	require.Equal(t, uint64(88), got.MostUsedTag.ID)
}

func TestGetUserStatsWithoutMostUsedEntities(t *testing.T) {
	repo, dbMock, _ := newUserRepo(t)

	expectCountQuery := func(rawArgs int) {
		dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
		if rawArgs == 1 {
			dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
		} else {
			dbMock.EXPECT().Raw(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbMock)
		}
		dbMock.EXPECT().Scan(gomock.Any()).Return(dbMock)
	}

	expectCountQuery(1)
	expectCountQuery(1)
	expectCountQuery(1)
	expectCountQuery(2)
	expectCountQuery(2)

	dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Scan(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Error().Return(errors.New("category query failed"))

	dbMock.EXPECT().WithContext(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Raw(gomock.Any(), gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Scan(gomock.Any()).Return(dbMock)
	dbMock.EXPECT().Error().Return(errors.New("tag query failed"))

	got, err := repo.GetUserStats(t.Context(), 100)
	require.NoError(t, err)
	require.Nil(t, got.MostUsedCategory)
	require.Nil(t, got.MostUsedTag)
}
