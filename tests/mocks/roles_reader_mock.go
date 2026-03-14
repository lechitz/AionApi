// Package mocks contains gomock mocks used across unit tests.
package mocks

import (
	"context"
	"reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockRolesReader is a mock of RolesReader interface.
//
// Note: this file is kept small and stable; most mocks are generated via mockgen.
type MockRolesReader struct {
	ctrl     *gomock.Controller
	recorder *MockRolesReaderMockRecorder
}

// MockRolesReaderMockRecorder is the mock recorder for MockRolesReader.
type MockRolesReaderMockRecorder struct {
	mock *MockRolesReader
}

// NewMockRolesReader creates a new mock instance.
func NewMockRolesReader(ctrl *gomock.Controller) *MockRolesReader {
	mock := &MockRolesReader{ctrl: ctrl}
	mock.recorder = &MockRolesReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRolesReader) EXPECT() *MockRolesReaderMockRecorder {
	return m.recorder
}

// GetRolesByUserID mocks base method.
func (m *MockRolesReader) GetRolesByUserID(ctx context.Context, userID uint64) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRolesByUserID", ctx, userID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRolesByUserID indicates an expected call of GetRolesByUserID.
func (mr *MockRolesReaderMockRecorder) GetRolesByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRolesByUserID", reflect.TypeOf((*MockRolesReader)(nil).GetRolesByUserID), ctx, userID)
}
