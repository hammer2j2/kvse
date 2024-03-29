// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hammer2j2/kvse/pkg/kvse (interfaces: Kvse)

// Package mock_kvse is a generated GoMock package.
package mock_kvse

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockKvse is a mock of Kvse interface.
type MockKvse struct {
	ctrl     *gomock.Controller
	recorder *MockKvseMockRecorder
}

// MockKvseMockRecorder is the mock recorder for MockKvse.
type MockKvseMockRecorder struct {
	mock *MockKvse
}

// NewMockKvse creates a new mock instance.
func NewMockKvse(ctrl *gomock.Controller) *MockKvse {
	mock := &MockKvse{ctrl: ctrl}
	mock.recorder = &MockKvseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKvse) EXPECT() *MockKvseMockRecorder {
	return m.recorder
}

// Read mocks base method.
func (m *MockKvse) Read() (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read")
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockKvseMockRecorder) Read() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockKvse)(nil).Read))
}

// SetupRequest mocks base method.
func (m *MockKvse) SetupRequest(arg0, arg1 string, arg2 []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetupRequest", arg0, arg1, arg2)
}

// SetupRequest indicates an expected call of SetupRequest.
func (mr *MockKvseMockRecorder) SetupRequest(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetupRequest", reflect.TypeOf((*MockKvse)(nil).SetupRequest), arg0, arg1, arg2)
}
