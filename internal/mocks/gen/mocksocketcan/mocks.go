// Code generated by MockGen. DO NOT EDIT.
// Source: ../../pkg/socketcan/fileconn.go

// Package mocksocketcan is a generated GoMock package.
package mocksocketcan

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// Mockfile is a mock of file interface.
type Mockfile struct {
	ctrl     *gomock.Controller
	recorder *MockfileMockRecorder
}

// MockfileMockRecorder is the mock recorder for Mockfile.
type MockfileMockRecorder struct {
	mock *Mockfile
}

// NewMockfile creates a new mock instance.
func NewMockfile(ctrl *gomock.Controller) *Mockfile {
	mock := &Mockfile{ctrl: ctrl}
	mock.recorder = &MockfileMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockfile) EXPECT() *MockfileMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *Mockfile) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockfileMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*Mockfile)(nil).Close))
}

// Read mocks base method.
func (m *Mockfile) Read(arg0 []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockfileMockRecorder) Read(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*Mockfile)(nil).Read), arg0)
}

// SetDeadline mocks base method.
func (m *Mockfile) SetDeadline(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetDeadline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDeadline indicates an expected call of SetDeadline.
func (mr *MockfileMockRecorder) SetDeadline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDeadline", reflect.TypeOf((*Mockfile)(nil).SetDeadline), arg0)
}

// SetReadDeadline mocks base method.
func (m *Mockfile) SetReadDeadline(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetReadDeadline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetReadDeadline indicates an expected call of SetReadDeadline.
func (mr *MockfileMockRecorder) SetReadDeadline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReadDeadline", reflect.TypeOf((*Mockfile)(nil).SetReadDeadline), arg0)
}

// SetWriteDeadline mocks base method.
func (m *Mockfile) SetWriteDeadline(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetWriteDeadline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetWriteDeadline indicates an expected call of SetWriteDeadline.
func (mr *MockfileMockRecorder) SetWriteDeadline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWriteDeadline", reflect.TypeOf((*Mockfile)(nil).SetWriteDeadline), arg0)
}

// Write mocks base method.
func (m *Mockfile) Write(arg0 []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockfileMockRecorder) Write(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*Mockfile)(nil).Write), arg0)
}
