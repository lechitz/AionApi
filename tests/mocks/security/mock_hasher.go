// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/ports/output/security/password.go

// Package security is a generated GoMock package.
package security

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHasher is a mock of Hasher interface.
type MockHasher struct {
	ctrl     *gomock.Controller
	recorder *MockHasherMockRecorder
}

// MockHasherMockRecorder is the mock recorder for MockHasher.
type MockHasherMockRecorder struct {
	mock *MockHasher
}

// NewMockHasher creates a new mock instance.
func NewMockHasher(ctrl *gomock.Controller) *MockHasher {
	mock := &MockHasher{ctrl: ctrl}
	mock.recorder = &MockHasherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHasher) EXPECT() *MockHasherMockRecorder {
	return m.recorder
}

// ComparePasswords mocks base method.
func (m *MockHasher) ComparePasswords(hashed, plain string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ComparePasswords", hashed, plain)
	ret0, _ := ret[0].(error)
	return ret0
}

// ComparePasswords indicates an expected call of ComparePasswords.
func (mr *MockHasherMockRecorder) ComparePasswords(hashed, plain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ComparePasswords", reflect.TypeOf((*MockHasher)(nil).ComparePasswords), hashed, plain)
}

// HashPassword mocks base method.
func (m *MockHasher) HashPassword(plain string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HashPassword", plain)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HashPassword indicates an expected call of HashPassword.
func (mr *MockHasherMockRecorder) HashPassword(plain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HashPassword", reflect.TypeOf((*MockHasher)(nil).HashPassword), plain)
}
