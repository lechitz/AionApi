// Package mocks contains gomock mocks used across unit tests.
package mocks

import (
	"context"
	"reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockSessionRevoker is a mock of SessionRevoker interface.
type MockSessionRevoker struct {
	ctrl     *gomock.Controller
	recorder *MockSessionRevokerMockRecorder
}

// MockSessionRevokerMockRecorder is the mock recorder for MockSessionRevoker.
type MockSessionRevokerMockRecorder struct {
	mock *MockSessionRevoker
}

// NewMockSessionRevoker creates a new mock instance.
func NewMockSessionRevoker(ctrl *gomock.Controller) *MockSessionRevoker {
	mock := &MockSessionRevoker{ctrl: ctrl}
	mock.recorder = &MockSessionRevokerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionRevoker) EXPECT() *MockSessionRevokerMockRecorder {
	return m.recorder
}

// RevokeUserSessions mocks base method.
func (m *MockSessionRevoker) RevokeUserSessions(ctx context.Context, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RevokeUserSessions", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RevokeUserSessions indicates an expected call of RevokeUserSessions.
func (mr *MockSessionRevokerMockRecorder) RevokeUserSessions(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokeUserSessions", reflect.TypeOf((*MockSessionRevoker)(nil).RevokeUserSessions), ctx, userID)
}
