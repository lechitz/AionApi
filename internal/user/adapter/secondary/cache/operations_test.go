package cache_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	usercache "github.com/lechitz/AionApi/internal/user/adapter/secondary/cache"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"github.com/stretchr/testify/require"
)

type fakeCache struct {
	data      map[string]string
	setErr    map[string]error
	getErr    map[string]error
	delErr    map[string]error
	lastTTL   map[string]time.Duration
	setCalled []string
	delCalled []string
}

func newFakeCache() *fakeCache {
	return &fakeCache{
		data:    make(map[string]string),
		setErr:  make(map[string]error),
		getErr:  make(map[string]error),
		delErr:  make(map[string]error),
		lastTTL: make(map[string]time.Duration),
	}
}

func (f *fakeCache) Set(_ context.Context, key string, value interface{}, exp time.Duration) error {
	if err := f.setErr[key]; err != nil {
		return err
	}
	f.setCalled = append(f.setCalled, key)
	f.lastTTL[key] = exp
	f.data[key] = fmt.Sprintf("%v", value)
	return nil
}

func (f *fakeCache) Get(_ context.Context, key string) (string, error) {
	if err := f.getErr[key]; err != nil {
		return "", err
	}
	v, ok := f.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

func (f *fakeCache) Del(_ context.Context, key string) error {
	if err := f.delErr[key]; err != nil {
		return err
	}
	f.delCalled = append(f.delCalled, key)
	delete(f.data, key)
	return nil
}

func (f *fakeCache) Ping(context.Context) error { return nil }
func (f *fakeCache) Close() error               { return nil }

func newUserFixture() domain.User {
	now := time.Now().UTC().Truncate(time.Second)
	return domain.User{
		ID:        7,
		Name:      "John Doe",
		Username:  "john",
		Email:     "john@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestSaveUser_UsesDefaultExpirationWhenNonPositive(t *testing.T) {
	fc := newFakeCache()
	store := usercache.NewStore(fc, &mockLogger{})
	user := newUserFixture()

	err := store.SaveUser(t.Context(), user, 0)
	require.NoError(t, err)

	idKey := fmt.Sprintf("user:id:%d", user.ID)
	require.Equal(t, usercache.UserExpirationDefault, fc.lastTTL[idKey])
}

func TestSaveUser_PrimaryKeyFailureReturnsError(t *testing.T) {
	fc := newFakeCache()
	user := newUserFixture()
	idKey := fmt.Sprintf("user:id:%d", user.ID)
	fc.setErr[idKey] = errors.New("redis down")

	store := usercache.NewStore(fc, &mockLogger{})
	err := store.SaveUser(t.Context(), user, time.Minute)
	require.Error(t, err)
}

func TestDeleteUser_PrimaryKeyFailureReturnsError(t *testing.T) {
	fc := newFakeCache()
	fc.delErr["user:id:11"] = errors.New("del failed")

	store := usercache.NewStore(fc, &mockLogger{})
	err := store.DeleteUser(t.Context(), 11, "john", "john@example.com")
	require.Error(t, err)
}

func TestDeleteUser_SecondaryKeyFailuresAreBestEffort(t *testing.T) {
	fc := newFakeCache()
	fc.delErr["user:username:john"] = errors.New("warn only")
	fc.delErr["user:email:john@example.com"] = errors.New("warn only")

	store := usercache.NewStore(fc, &mockLogger{})
	err := store.DeleteUser(t.Context(), 11, "john", "john@example.com")
	require.NoError(t, err)
}

func TestGetUserByUsername_AndByEmail_Success(t *testing.T) {
	fc := newFakeCache()
	store := usercache.NewStore(fc, &mockLogger{})
	user := newUserFixture()
	require.NoError(t, store.SaveUser(t.Context(), user, time.Minute))

	byUsername, err := store.GetUserByUsername(t.Context(), user.Username)
	require.NoError(t, err)
	require.Equal(t, user.ID, byUsername.ID)
	require.Equal(t, user.Email, byUsername.Email)

	byEmail, err := store.GetUserByEmail(t.Context(), user.Email)
	require.NoError(t, err)
	require.Equal(t, user.ID, byEmail.ID)
	require.Equal(t, user.Username, byEmail.Username)
}

func TestGetUserByUsername_InvalidJSONReturnsError(t *testing.T) {
	fc := newFakeCache()
	fc.data["user:username:john"] = "{invalid-json"

	store := usercache.NewStore(fc, &mockLogger{})
	_, err := store.GetUserByUsername(t.Context(), "john")
	require.Error(t, err)
}

func TestGetUserByEmail_StaleVersionReturnsError(t *testing.T) {
	fc := newFakeCache()
	fc.data["user:email:john@example.com"] = `{"version":999,"id":7,"username":"john","email":"john@example.com"}`

	store := usercache.NewStore(fc, &mockLogger{})
	_, err := store.GetUserByEmail(t.Context(), "john@example.com")
	require.Error(t, err)
}

func TestGetUserByID_CacheAndPayloadErrors(t *testing.T) {
	fc := newFakeCache()
	store := usercache.NewStore(fc, &mockLogger{})

	idKey := "user:id:77"
	fc.getErr[idKey] = errors.New("redis get failed")
	_, err := store.GetUserByID(t.Context(), 77)
	require.Error(t, err)

	delete(fc.getErr, idKey)
	fc.data[idKey] = "{invalid-json"
	_, err = store.GetUserByID(t.Context(), 77)
	require.Error(t, err)

	fc.data[idKey] = `{"version":999,"id":77,"username":"john","email":"john@example.com"}`
	_, err = store.GetUserByID(t.Context(), 77)
	require.Error(t, err)
}

func TestGetUserByUsernameAndEmail_CacheGetErrors(t *testing.T) {
	fc := newFakeCache()
	store := usercache.NewStore(fc, &mockLogger{})

	fc.getErr["user:username:john"] = errors.New("username get failed")
	_, err := store.GetUserByUsername(t.Context(), "john")
	require.Error(t, err)

	fc.getErr["user:email:john@example.com"] = errors.New("email get failed")
	_, err = store.GetUserByEmail(t.Context(), "john@example.com")
	require.Error(t, err)
}
