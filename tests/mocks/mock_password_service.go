// Code generated by MockGen. DO NOT EDIT.
// Source: ports/output/security/password.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIPasswordService is a mock of IPasswordService interface.
type MockIPasswordService struct {
	ctrl     *gomock.Controller
	recorder *MockIPasswordServiceMockRecorder
}

// MockIPasswordServiceMockRecorder is the mock recorder for MockIPasswordService.
type MockIPasswordServiceMockRecorder struct {
	mock *MockIPasswordService
}

// NewMockIPasswordService creates a new mock instance.
func NewMockIPasswordService(ctrl *gomock.Controller) *MockIPasswordService {
	mock := &MockIPasswordService{ctrl: ctrl}
	mock.recorder = &MockIPasswordServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPasswordService) EXPECT() *MockIPasswordServiceMockRecorder {
	return m.recorder
}

// ComparePasswords mocks base method.
func (m *MockIPasswordService) ComparePasswords(hashed, plain string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ComparePasswords", hashed, plain)
	ret0, _ := ret[0].(error)
	return ret0
}

// ComparePasswords indicates an expected call of ComparePasswords.
func (mr *MockIPasswordServiceMockRecorder) ComparePasswords(hashed, plain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ComparePasswords", reflect.TypeOf((*MockIPasswordService)(nil).ComparePasswords), hashed, plain)
}

// HashPassword mocks base method.
func (m *MockIPasswordService) HashPassword(plain string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HashPassword", plain)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HashPassword indicates an expected call of HashPassword.
func (mr *MockIPasswordServiceMockRecorder) HashPassword(plain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HashPassword", reflect.TypeOf((*MockIPasswordService)(nil).HashPassword), plain)
}
