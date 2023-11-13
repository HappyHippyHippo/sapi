package sapi

import (
	"reflect"

	"github.com/golang/mock/gomock"
	"github.com/happyhippyhippo/slate"
)

// ----------------------------------------------------------------------------
// Config Supplier
// ----------------------------------------------------------------------------

// MockConfigSupplier is a mock instance of config.Source interface
type MockConfigSupplier struct {
	ctrl     *gomock.Controller
	recorder *MockConfigSupplierRecorder
}

var _ slate.ConfigSupplier = &MockConfigSupplier{}

// MockConfigSupplierRecorder is the mock recorder for MockConfigSupplier
type MockConfigSupplierRecorder struct {
	mock *MockConfigSupplier
}

// NewMockConfigSupplier creates a new mock instance
func NewMockConfigSupplier(ctrl *gomock.Controller) *MockConfigSupplier {
	mock := &MockConfigSupplier{ctrl: ctrl}
	mock.recorder = &MockConfigSupplierRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfigSupplier) EXPECT() *MockConfigSupplierRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockConfigSupplier) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockConfigSupplierRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConfigSupplier)(nil).Close))
}

// Has mocks base method
func (m *MockConfigSupplier) Has(path string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Has", path)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Has indicates an expected call of Has
func (mr *MockConfigSupplierRecorder) Has(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Has", reflect.TypeOf((*MockConfigSupplier)(nil).Has), path)
}

// Get mocks base method
func (m *MockConfigSupplier) Get(path string, def ...interface{}) (interface{}, error) {
	m.ctrl.T.Helper()
	var varargs []interface{}
	varargs = append(varargs, path)
	for _, a := range def {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockConfigSupplierRecorder) Get(path interface{}, def ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	var varargs []interface{}
	varargs = append(varargs, path)
	for _, a := range def {
		varargs = append(varargs, a)
	}
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockConfigSupplier)(nil).Get), varargs...)
}

// ----------------------------------------------------------------------------
// Log Writer
// ----------------------------------------------------------------------------

// MockLogWriter is a mock instance of Stream interface.
type MockLogWriter struct {
	ctrl     *gomock.Controller
	recorder *MockLogWriterRecorder
}

var _ slate.LogWriter = &MockLogWriter{}

// MockLogWriterRecorder is the mock recorder for MockLogWriter.
type MockLogWriterRecorder struct {
	mock *MockLogWriter
}

// NewMockLogWriter creates a new mock instance.
func NewMockLogWriter(ctrl *gomock.Controller) *MockLogWriter {
	mock := &MockLogWriter{ctrl: ctrl}
	mock.recorder = &MockLogWriterRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogWriter) EXPECT() *MockLogWriterRecorder {
	return m.recorder
}

// AddChannel mocks base method.
func (m *MockLogWriter) AddChannel(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddChannel", arg0)
}

// AddChannel indicates an expected call of AddChannel.
func (mr *MockLogWriterRecorder) AddChannel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddChannel", reflect.TypeOf((*MockLogWriter)(nil).AddChannel), arg0)
}

// Broadcast mocks base method.
func (m *MockLogWriter) Broadcast(arg0 slate.LogLevel, arg1 string, arg2 ...slate.LogContext) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Broadcast", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Broadcast indicates an expected call of Broadcast.
func (mr *MockLogWriterRecorder) Broadcast(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Broadcast", reflect.TypeOf((*MockLogWriter)(nil).Broadcast), varargs...)
}

// Flush mocks base method.
func (m *MockLogWriter) Flush() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Flush")
	ret0, _ := ret[0].(error)
	return ret0
}

// Flush indicates an expected call of Flush.
func (mr *MockLogWriterRecorder) Flush() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flush", reflect.TypeOf((*MockLogWriter)(nil).Flush))
}

// HasChannel mocks base method.
func (m *MockLogWriter) HasChannel(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasChannel", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasChannel indicates an expected call of HasChannel.
func (mr *MockLogWriterRecorder) HasChannel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasChannel", reflect.TypeOf((*MockLogWriter)(nil).HasChannel), arg0)
}

// ListChannels mocks base method.
func (m *MockLogWriter) ListChannels() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListChannels")
	ret0, _ := ret[0].([]string)
	return ret0
}

// ListChannels indicates an expected call of ListChannels.
func (mr *MockLogWriterRecorder) ListChannels() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListChannels", reflect.TypeOf((*MockLogWriter)(nil).ListChannels))
}

// RemoveChannel mocks base method.
func (m *MockLogWriter) RemoveChannel(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveChannel", arg0)
}

// RemoveChannel indicates an expected call of RemoveChannel.
func (mr *MockLogWriterRecorder) RemoveChannel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveChannel", reflect.TypeOf((*MockLogWriter)(nil).RemoveChannel), arg0)
}

// Signal mocks base method.
func (m *MockLogWriter) Signal(arg0 string, arg1 slate.LogLevel, arg2 string, arg3 ...slate.LogContext) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Signal", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Signal indicates an expected call of Signal.
func (mr *MockLogWriterRecorder) Signal(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signal", reflect.TypeOf((*MockLogWriter)(nil).Signal), varargs...)
}
