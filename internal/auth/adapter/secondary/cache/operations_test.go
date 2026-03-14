package cache_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	authcache "github.com/lechitz/AionApi/internal/auth/adapter/secondary/cache"
	"github.com/lechitz/AionApi/internal/auth/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/stretchr/testify/require"
)

type fakeAuthCache struct {
	data      map[string]string
	setErr    map[string]error
	getErr    map[string]error
	delErr    map[string]error
	lastTTL   map[string]time.Duration
	setCalled []string
	delCalled []string
}

func newFakeAuthCache() *fakeAuthCache {
	return &fakeAuthCache{
		data:    make(map[string]string),
		setErr:  make(map[string]error),
		getErr:  make(map[string]error),
		delErr:  make(map[string]error),
		lastTTL: make(map[string]time.Duration),
	}
}

func (f *fakeAuthCache) Set(_ context.Context, key string, value interface{}, exp time.Duration) error {
	if err := f.setErr[key]; err != nil {
		return err
	}
	f.setCalled = append(f.setCalled, key)
	f.lastTTL[key] = exp
	f.data[key] = fmt.Sprintf("%v", value)
	return nil
}

func (f *fakeAuthCache) Get(_ context.Context, key string) (string, error) {
	if err := f.getErr[key]; err != nil {
		return "", err
	}
	return f.data[key], nil
}

func (f *fakeAuthCache) Del(_ context.Context, key string) error {
	if err := f.delErr[key]; err != nil {
		return err
	}
	f.delCalled = append(f.delCalled, key)
	delete(f.data, key)
	return nil
}

func (f *fakeAuthCache) Ping(context.Context) error { return nil }
func (f *fakeAuthCache) Close() error               { return nil }

type mockAuthLogger struct{}

func (m *mockAuthLogger) Infow(_ string, _ ...interface{})                        {}
func (m *mockAuthLogger) Warnw(_ string, _ ...interface{})                        {}
func (m *mockAuthLogger) Errorw(_ string, _ ...interface{})                       {}
func (m *mockAuthLogger) InfowCtx(_ context.Context, _ string, _ ...interface{})  {}
func (m *mockAuthLogger) WarnwCtx(_ context.Context, _ string, _ ...interface{})  {}
func (m *mockAuthLogger) ErrorwCtx(_ context.Context, _ string, _ ...interface{}) {}
func (m *mockAuthLogger) DebugwCtx(_ context.Context, _ string, _ ...interface{}) {}
func (m *mockAuthLogger) Debugw(_ string, _ ...interface{})                       {}
func (m *mockAuthLogger) Debugf(_ string, _ ...interface{})                       {}
func (m *mockAuthLogger) Infof(_ string, _ ...interface{})                        {}
func (m *mockAuthLogger) Warnf(_ string, _ ...interface{})                        {}
func (m *mockAuthLogger) Errorf(_ string, _ ...interface{})                       {}
func (m *mockAuthLogger) Sync() error                                             { return nil }

func TestAuthCache_SaveAndGetAccessToken(t *testing.T) {
	fc := newFakeAuthCache()
	store := authcache.NewStore(fc, &mockAuthLogger{})
	token := domain.NewAccessToken("jwt-access", 42).ToAuth()

	err := store.Save(t.Context(), token, 0)
	require.NoError(t, err)

	key := fmt.Sprintf(authcache.TokenUserKeyFormat, token.Key, token.Type)
	require.Equal(t, authcache.TokenExpirationDefault, fc.lastTTL[key])

	got, err := store.Get(t.Context(), token.Key, commonkeys.TokenTypeAccess)
	require.NoError(t, err)
	require.Equal(t, token.Key, got.Key)
	require.Equal(t, token.Token, got.Token)
	require.Equal(t, domain.TokenTypeAccess, got.Type)
}

func TestAuthCache_GetRefreshAndUnknownType(t *testing.T) {
	fc := newFakeAuthCache()
	store := authcache.NewStore(fc, &mockAuthLogger{})

	refreshKey := fmt.Sprintf(authcache.TokenUserKeyFormat, 9, commonkeys.TokenTypeRefresh)
	fc.data[refreshKey] = "refresh-token"

	gotRefresh, err := store.Get(t.Context(), 9, commonkeys.TokenTypeRefresh)
	require.NoError(t, err)
	require.Equal(t, "refresh-token", gotRefresh.Token)
	require.Equal(t, domain.TokenTypeRefresh, gotRefresh.Type)

	unknownKey := fmt.Sprintf(authcache.TokenUserKeyFormat, 9, "custom")
	fc.data[unknownKey] = "custom-token"
	gotUnknown, err := store.Get(t.Context(), 9, "custom")
	require.NoError(t, err)
	require.Equal(t, "custom-token", gotUnknown.Token)
	require.Equal(t, "custom", gotUnknown.Type)
}

func TestAuthCache_GetReturnsErrorOnCacheFailure(t *testing.T) {
	fc := newFakeAuthCache()
	key := fmt.Sprintf(authcache.TokenUserKeyFormat, 7, commonkeys.TokenTypeAccess)
	fc.getErr[key] = errors.New("redis get failed")
	store := authcache.NewStore(fc, &mockAuthLogger{})

	_, err := store.Get(t.Context(), 7, commonkeys.TokenTypeAccess)
	require.Error(t, err)
}

func TestAuthCache_DeleteAndRevokeSessions(t *testing.T) {
	fc := newFakeAuthCache()
	store := authcache.NewStore(fc, &mockAuthLogger{})

	require.NoError(t, store.Delete(t.Context(), 11, commonkeys.TokenTypeAccess))
	require.NoError(t, store.RevokeUserSessions(t.Context(), 11))

	accessKey := fmt.Sprintf(authcache.TokenUserKeyFormat, 11, commonkeys.TokenTypeAccess)
	refreshKey := fmt.Sprintf(authcache.TokenUserKeyFormat, 11, commonkeys.TokenTypeRefresh)
	require.Contains(t, fc.delCalled, accessKey)
	require.Contains(t, fc.delCalled, refreshKey)
}

func TestAuthCache_RevokeSessionsStopsOnFirstError(t *testing.T) {
	fc := newFakeAuthCache()
	accessKey := fmt.Sprintf(authcache.TokenUserKeyFormat, 11, commonkeys.TokenTypeAccess)
	fc.delErr[accessKey] = errors.New("delete access failed")
	store := authcache.NewStore(fc, &mockAuthLogger{})

	err := store.RevokeUserSessions(t.Context(), 11)
	require.Error(t, err)
}

func TestAuthCache_SaveWithKeyAndGetByKey(t *testing.T) {
	fc := newFakeAuthCache()
	store := authcache.NewStore(fc, &mockAuthLogger{})

	customKey := "token:grace:user:44:access"
	err := store.SaveWithKey(t.Context(), customKey, domain.Auth{Token: "grace-token"}, 0)
	require.NoError(t, err)
	require.Equal(t, authcache.TokenExpirationDefault, fc.lastTTL[customKey])

	got, err := store.GetByKey(t.Context(), customKey)
	require.NoError(t, err)
	require.Equal(t, "grace-token", got.Token)
}

func TestAuthCache_GetByKeyNotFoundAndError(t *testing.T) {
	fc := newFakeAuthCache()
	store := authcache.NewStore(fc, &mockAuthLogger{})

	_, err := store.GetByKey(t.Context(), "missing:key")
	require.Error(t, err)

	fc.getErr["broken:key"] = errors.New("redis unavailable")
	_, err = store.GetByKey(t.Context(), "broken:key")
	require.Error(t, err)
}

func TestAuthCache_RoleOperations(t *testing.T) {
	fc := newFakeAuthCache()
	store := authcache.NewStore(fc, &mockAuthLogger{})
	userID := uint64(77)
	roleKey := fmt.Sprintf(authcache.RoleKeyFormat, userID)

	require.NoError(t, store.SaveRoles(t.Context(), userID, []string{"admin", "reader"}, 0))
	require.Equal(t, authcache.RoleExpirationDefault, fc.lastTTL[roleKey])

	roles, err := store.GetRoles(t.Context(), userID)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"admin", "reader"}, roles)

	require.NoError(t, store.InvalidateRoles(t.Context(), userID))
	require.Contains(t, fc.delCalled, roleKey)
}

func TestAuthCache_RoleGetAndSaveErrors(t *testing.T) {
	fc := newFakeAuthCache()
	store := authcache.NewStore(fc, &mockAuthLogger{})
	userID := uint64(88)
	roleKey := fmt.Sprintf(authcache.RoleKeyFormat, userID)

	fc.setErr[roleKey] = errors.New("set roles failed")
	err := store.SaveRoles(t.Context(), userID, []string{"admin"}, time.Minute)
	require.Error(t, err)
	delete(fc.setErr, roleKey)

	fc.data[roleKey] = "not-json"
	_, err = store.GetRoles(t.Context(), userID)
	require.Error(t, err)
	delete(fc.data, roleKey)

	fc.getErr[roleKey] = errors.New("get roles failed")
	_, err = store.GetRoles(t.Context(), userID)
	require.Error(t, err)

	fc.delErr[roleKey] = errors.New("del roles failed")
	err = store.InvalidateRoles(t.Context(), userID)
	require.Error(t, err)
}
