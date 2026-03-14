// Package mocks contains gomock mocks used across unit tests.
package mocks

import (
	"context"
	"reflect"
	"time"

	gomock "go.uber.org/mock/gomock"
)

// MockRoleCache is a mock of RoleCache interface.
type MockRoleCache struct {
	ctrl     *gomock.Controller
	recorder *MockRoleCacheMockRecorder
}

// MockRoleCacheMockRecorder is the mock recorder for MockRoleCache.
type MockRoleCacheMockRecorder struct {
	mock *MockRoleCache
}

// NewMockRoleCache creates a new mock instance.
func NewMockRoleCache(ctrl *gomock.Controller) *MockRoleCache {
	mock := &MockRoleCache{ctrl: ctrl}
	mock.recorder = &MockRoleCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRoleCache) EXPECT() *MockRoleCacheMockRecorder {
	return m.recorder
}

// SaveRoles mocks base method.
func (m *MockRoleCache) SaveRoles(ctx context.Context, userID uint64, roles []string, ttl time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveRoles", ctx, userID, roles, ttl)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveRoles indicates an expected call of SaveRoles.
func (mr *MockRoleCacheMockRecorder) SaveRoles(ctx, userID, roles, ttl any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveRoles", reflect.TypeOf((*MockRoleCache)(nil).SaveRoles), ctx, userID, roles, ttl)
}

// GetRoles mocks base method.
func (m *MockRoleCache) GetRoles(ctx context.Context, userID uint64) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRoles", ctx, userID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRoles indicates an expected call of GetRoles.
func (mr *MockRoleCacheMockRecorder) GetRoles(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoles", reflect.TypeOf((*MockRoleCache)(nil).GetRoles), ctx, userID)
}

// InvalidateRoles mocks base method.
func (m *MockRoleCache) InvalidateRoles(ctx context.Context, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvalidateRoles", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// InvalidateRoles indicates an expected call of InvalidateRoles.
func (mr *MockRoleCacheMockRecorder) InvalidateRoles(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvalidateRoles", reflect.TypeOf((*MockRoleCache)(nil).InvalidateRoles), ctx, userID)
}
