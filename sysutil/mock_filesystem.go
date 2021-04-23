// Code generated by MockGen. DO NOT EDIT.
// Source: filesystem.go

// Package mock_sysutil is a generated GoMock package.
package sysutil

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFileSystemInterface is a mock of FileSystemInterface interface.
type MockFileSystemInterface struct {
	ctrl     *gomock.Controller
	recorder *MockFileSystemInterfaceMockRecorder
}

// MockFileSystemInterfaceMockRecorder is the mock recorder for MockFileSystemInterface.
type MockFileSystemInterfaceMockRecorder struct {
	mock *MockFileSystemInterface
}

// NewMockFileSystemInterface creates a new mock instance.
func NewMockFileSystemInterface(ctrl *gomock.Controller) *MockFileSystemInterface {
	mock := &MockFileSystemInterface{ctrl: ctrl}
	mock.recorder = &MockFileSystemInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileSystemInterface) EXPECT() *MockFileSystemInterfaceMockRecorder {
	return m.recorder
}

// ListFolders mocks base method.
func (m *MockFileSystemInterface) ListFolders(path string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFolders", path)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFolders indicates an expected call of ListFolders.
func (mr *MockFileSystemInterfaceMockRecorder) ListFolders(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFolders", reflect.TypeOf((*MockFileSystemInterface)(nil).ListFolders), path)
}
