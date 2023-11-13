package sapi

import (
	"html/template"
	"net"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

// ----------------------------------------------------------------------------
// RestEngine
// ----------------------------------------------------------------------------

// MockRestEngine is a mock instance of Engine interface.
type MockRestEngine struct {
	ctrl     *gomock.Controller
	recorder *MockRestEngineRecorder
}

var _ RestEngine = &MockRestEngine{}

// MockRestEngineRecorder is the mock recorder for MockRestEngine.
type MockRestEngineRecorder struct {
	mock *MockRestEngine
}

// NewMockRestEngine creates a new mock instance.
func NewMockRestEngine(ctrl *gomock.Controller) *MockRestEngine {
	mock := &MockRestEngine{ctrl: ctrl}
	mock.recorder = &MockRestEngineRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRestEngine) EXPECT() *MockRestEngineRecorder {
	return m.recorder
}

// Any mocks base method.
func (m *MockRestEngine) Any(arg0 string, arg1 ...gin.HandlerFunc) gin.IRoutes {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Any", varargs...)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// Any indicates an expected call of Any.
func (mr *MockRestEngineRecorder) Any(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Any", reflect.TypeOf((*MockRestEngine)(nil).Any), varargs...)
}

// DELETE mocks base method.
func (m *MockRestEngine) DELETE(arg0 string, arg1 ...gin.HandlerFunc) gin.IRoutes {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DELETE", varargs...)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// DELETE indicates an expected call of DELETE.
func (mr *MockRestEngineRecorder) DELETE(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DELETE", reflect.TypeOf((*MockRestEngine)(nil).DELETE), varargs...)
}

// Delims mocks base method.
func (m *MockRestEngine) Delims(left, right string) *gin.Engine {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delims", left, right)
	ret0, _ := ret[0].(*gin.Engine)
	return ret0
}

// Delims indicates an expected call of Delims.
func (mr *MockRestEngineRecorder) Delims(left, right interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delims", reflect.TypeOf((*MockRestEngine)(nil).Delims), left, right)
}

// GET mocks base method.
func (m *MockRestEngine) GET(arg0 string, arg1 ...gin.HandlerFunc) gin.IRoutes {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GET", varargs...)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// GET indicates an expected call of GET.
func (mr *MockRestEngineRecorder) GET(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GET", reflect.TypeOf((*MockRestEngine)(nil).GET), varargs...)
}

// Group mocks base method.
func (m *MockRestEngine) Group(arg0 string, arg1 ...gin.HandlerFunc) *gin.RouterGroup {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Group", varargs...)
	ret0, _ := ret[0].(*gin.RouterGroup)
	return ret0
}

// Group indicates an expected call of Group.
func (mr *MockRestEngineRecorder) Group(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Group", reflect.TypeOf((*MockRestEngine)(nil).Group), varargs...)
}

// HEAD mocks base method.
func (m *MockRestEngine) HEAD(arg0 string, arg1 ...gin.HandlerFunc) gin.IRoutes {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "HEAD", varargs...)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// HEAD indicates an expected call of HEAD.
func (mr *MockRestEngineRecorder) HEAD(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HEAD", reflect.TypeOf((*MockRestEngine)(nil).HEAD), varargs...)
}

// Handle mocks base method.
func (m *MockRestEngine) Handle(arg0, arg1 string, arg2 ...gin.HandlerFunc) gin.IRoutes {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Handle", varargs...)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// Handle indicates an expected call of Handle.
func (mr *MockRestEngineRecorder) Handle(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockRestEngine)(nil).Handle), varargs...)
}

// HandleContext mocks base method.
func (m *MockRestEngine) HandleContext(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleContext", c)
}

// HandleContext indicates an expected call of HandleContext.
func (mr *MockRestEngineRecorder) HandleContext(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleContext", reflect.TypeOf((*MockRestEngine)(nil).HandleContext), c)
}

// Handler mocks base method.
func (m *MockRestEngine) Handler() http.Handler {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handler")
	ret0, _ := ret[0].(http.Handler)
	return ret0
}

// Handler indicates an expected call of Handler.
func (mr *MockRestEngineRecorder) Handler() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handler", reflect.TypeOf((*MockRestEngine)(nil).Handler))
}

// LoadHTMLFiles mocks base method.
func (m *MockRestEngine) LoadHTMLFiles(files ...string) {
	m.ctrl.T.Helper()
	var varargs []interface{}
	for _, a := range files {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "LoadHTMLFiles", varargs...)
}

// LoadHTMLFiles indicates an expected call of LoadHTMLFiles.
func (mr *MockRestEngineRecorder) LoadHTMLFiles(files ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadHTMLFiles", reflect.TypeOf((*MockRestEngine)(nil).LoadHTMLFiles), files...)
}

// LoadHTMLGlob mocks base method.
func (m *MockRestEngine) LoadHTMLGlob(pattern string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LoadHTMLGlob", pattern)
}

// LoadHTMLGlob indicates an expected call of LoadHTMLGlob.
func (mr *MockRestEngineRecorder) LoadHTMLGlob(pattern interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadHTMLGlob", reflect.TypeOf((*MockRestEngine)(nil).LoadHTMLGlob), pattern)
}

// Match mocks base method.
func (m *MockRestEngine) Match(arg0 []string, arg1 string, arg2 ...gin.HandlerFunc) gin.IRoutes {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Match", varargs...)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// Match indicates an expected call of Match.
func (mr *MockRestEngineRecorder) Match(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Match", reflect.TypeOf((*MockRestEngine)(nil).Match), varargs...)
}

// NoMethod mocks base method.
func (m *MockRestEngine) NoMethod(handlers ...gin.HandlerFunc) {
	m.ctrl.T.Helper()
	var varargs []interface{}
	for _, a := range handlers {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "NoMethod", varargs...)
}

// NoMethod indicates an expected call of NoMethod.
func (mr *MockRestEngineRecorder) NoMethod(handlers ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NoMethod", reflect.TypeOf((*MockRestEngine)(nil).NoMethod), handlers...)
}

// NoRoute mocks base method.
func (m *MockRestEngine) NoRoute(handlers ...gin.HandlerFunc) {
	m.ctrl.T.Helper()
	var varargs []interface{}
	for _, a := range handlers {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "NoRoute", varargs...)
}

// NoRoute indicates an expected call of NoRoute.
func (mr *MockRestEngineRecorder) NoRoute(handlers ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NoRoute", reflect.TypeOf((*MockRestEngine)(nil).NoRoute), handlers...)
}

// OPTIONS mocks base method.
func (m *MockRestEngine) OPTIONS(arg0 string, arg1 ...gin.HandlerFunc) gin.IRoutes {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "OPTIONS", varargs...)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// OPTIONS indicates an expected call of OPTIONS.
func (mr *MockRestEngineRecorder) OPTIONS(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OPTIONS", reflect.TypeOf((*MockRestEngine)(nil).OPTIONS), varargs...)
}

// PATCH mocks base method.
func (m *MockRestEngine) PATCH(arg0 string, arg1 ...gin.HandlerFunc) gin.IRoutes {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PATCH", varargs...)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// PATCH indicates an expected call of PATCH.
func (mr *MockRestEngineRecorder) PATCH(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PATCH", reflect.TypeOf((*MockRestEngine)(nil).PATCH), varargs...)
}

// POST mocks base method.
func (m *MockRestEngine) POST(arg0 string, arg1 ...gin.HandlerFunc) gin.IRoutes {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "POST", varargs...)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// POST indicates an expected call of POST.
func (mr *MockRestEngineRecorder) POST(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "POST", reflect.TypeOf((*MockRestEngine)(nil).POST), varargs...)
}

// PUT mocks base method.
func (m *MockRestEngine) PUT(arg0 string, arg1 ...gin.HandlerFunc) gin.IRoutes {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PUT", varargs...)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// PUT indicates an expected call of PUT.
func (mr *MockRestEngineRecorder) PUT(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PUT", reflect.TypeOf((*MockRestEngine)(nil).PUT), varargs...)
}

// Routes mocks base method.
func (m *MockRestEngine) Routes() gin.RoutesInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Routes")
	ret0, _ := ret[0].(gin.RoutesInfo)
	return ret0
}

// Routes indicates an expected call of Routes.
func (mr *MockRestEngineRecorder) Routes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Routes", reflect.TypeOf((*MockRestEngine)(nil).Routes))
}

// Run mocks base method.
func (m *MockRestEngine) Run(addr ...string) error {
	m.ctrl.T.Helper()
	var varargs []interface{}
	for _, a := range addr {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Run", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockRestEngineRecorder) Run(addr ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockRestEngine)(nil).Run), addr...)
}

// RunFd mocks base method.
func (m *MockRestEngine) RunFd(fd int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunFd", fd)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunFd indicates an expected call of RunFd.
func (mr *MockRestEngineRecorder) RunFd(fd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunFd", reflect.TypeOf((*MockRestEngine)(nil).RunFd), fd)
}

// RunListener mocks base method.
func (m *MockRestEngine) RunListener(listener net.Listener) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunListener", listener)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunListener indicates an expected call of RunListener.
func (mr *MockRestEngineRecorder) RunListener(listener interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunListener", reflect.TypeOf((*MockRestEngine)(nil).RunListener), listener)
}

// RunTLS mocks base method.
func (m *MockRestEngine) RunTLS(addr, certFile, keyFile string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunTLS", addr, certFile, keyFile)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunTLS indicates an expected call of RunTLS.
func (mr *MockRestEngineRecorder) RunTLS(addr, certFile, keyFile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunTLS", reflect.TypeOf((*MockRestEngine)(nil).RunTLS), addr, certFile, keyFile)
}

// RunUnix mocks base method.
func (m *MockRestEngine) RunUnix(file string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunUnix", file)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunUnix indicates an expected call of RunUnix.
func (mr *MockRestEngineRecorder) RunUnix(file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunUnix", reflect.TypeOf((*MockRestEngine)(nil).RunUnix), file)
}

// SecureJsonPrefix mocks base method.
//
//revive:disable
func (m *MockRestEngine) SecureJsonPrefix(prefix string) *gin.Engine {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SecureJsonPrefix", prefix)
	ret0, _ := ret[0].(*gin.Engine)
	return ret0
}

// SecureJsonPrefix indicates an expected call of SecureJsonPrefix.
func (mr *MockRestEngineRecorder) SecureJsonPrefix(prefix interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SecureJsonPrefix", reflect.TypeOf((*MockRestEngine)(nil).SecureJsonPrefix), prefix)
}

//revive:enable

// ServeHTTP mocks base method.
func (m *MockRestEngine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ServeHTTP", w, req)
}

// ServeHTTP indicates an expected call of ServeHTTP.
func (mr *MockRestEngineRecorder) ServeHTTP(w, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServeHTTP", reflect.TypeOf((*MockRestEngine)(nil).ServeHTTP), w, req)
}

// SetFuncMap mocks base method.
func (m *MockRestEngine) SetFuncMap(funcMap template.FuncMap) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetFuncMap", funcMap)
}

// SetFuncMap indicates an expected call of SetFuncMap.
func (mr *MockRestEngineRecorder) SetFuncMap(funcMap interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFuncMap", reflect.TypeOf((*MockRestEngine)(nil).SetFuncMap), funcMap)
}

// SetHTMLTemplate mocks base method.
func (m *MockRestEngine) SetHTMLTemplate(templ *template.Template) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetHTMLTemplate", templ)
}

// SetHTMLTemplate indicates an expected call of SetHTMLTemplate.
func (mr *MockRestEngineRecorder) SetHTMLTemplate(templ interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHTMLTemplate", reflect.TypeOf((*MockRestEngine)(nil).SetHTMLTemplate), templ)
}

// SetTrustedProxies mocks base method.
func (m *MockRestEngine) SetTrustedProxies(trustedProxies []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetTrustedProxies", trustedProxies)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetTrustedProxies indicates an expected call of SetTrustedProxies.
func (mr *MockRestEngineRecorder) SetTrustedProxies(trustedProxies interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrustedProxies", reflect.TypeOf((*MockRestEngine)(nil).SetTrustedProxies), trustedProxies)
}

// Static mocks base method.
func (m *MockRestEngine) Static(arg0, arg1 string) gin.IRoutes {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Static", arg0, arg1)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// Static indicates an expected call of Static.
func (mr *MockRestEngineRecorder) Static(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Static", reflect.TypeOf((*MockRestEngine)(nil).Static), arg0, arg1)
}

// StaticFS mocks base method.
func (m *MockRestEngine) StaticFS(arg0 string, arg1 http.FileSystem) gin.IRoutes {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StaticFS", arg0, arg1)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// StaticFS indicates an expected call of StaticFS.
func (mr *MockRestEngineRecorder) StaticFS(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StaticFS", reflect.TypeOf((*MockRestEngine)(nil).StaticFS), arg0, arg1)
}

// StaticFile mocks base method.
func (m *MockRestEngine) StaticFile(arg0, arg1 string) gin.IRoutes {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StaticFile", arg0, arg1)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// StaticFile indicates an expected call of StaticFile.
func (mr *MockRestEngineRecorder) StaticFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StaticFile", reflect.TypeOf((*MockRestEngine)(nil).StaticFile), arg0, arg1)
}

// StaticFileFS mocks base method.
func (m *MockRestEngine) StaticFileFS(arg0, arg1 string, arg2 http.FileSystem) gin.IRoutes {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StaticFileFS", arg0, arg1, arg2)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// StaticFileFS indicates an expected call of StaticFileFS.
func (mr *MockRestEngineRecorder) StaticFileFS(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StaticFileFS", reflect.TypeOf((*MockRestEngine)(nil).StaticFileFS), arg0, arg1, arg2)
}

// Use mocks base method.
func (m *MockRestEngine) Use(arg0 ...gin.HandlerFunc) gin.IRoutes {
	m.ctrl.T.Helper()
	var varargs []interface{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Use", varargs...)
	ret0, _ := ret[0].(gin.IRoutes)
	return ret0
}

// Use indicates an expected call of Use.
func (mr *MockRestEngineRecorder) Use(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Use", reflect.TypeOf((*MockRestEngine)(nil).Use), arg0...)
}

// ----------------------------------------------------------------------------
// Rest Endpoint Registers
// ----------------------------------------------------------------------------

// MockRestEndpointRegister is a mock instance of Engine interface.
type MockRestEndpointRegister struct {
	ctrl     *gomock.Controller
	recorder *MockRestEndpointRegisterRecorder
}

var _ RestEndpointRegister = &MockRestEndpointRegister{}

// MockRestEndpointRegisterRecorder is the mock recorder for MockRestEngine.
type MockRestEndpointRegisterRecorder struct {
	mock *MockRestEndpointRegister
}

// NewMockRestEndpointRegister creates a new mock instance.
func NewMockRestEndpointRegister(ctrl *gomock.Controller) *MockRestEndpointRegister {
	mock := &MockRestEndpointRegister{ctrl: ctrl}
	mock.recorder = &MockRestEndpointRegisterRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRestEndpointRegister) EXPECT() *MockRestEndpointRegisterRecorder {
	return m.recorder
}

// Reg mocks base method.
func (m *MockRestEndpointRegister) Reg(arg0 RestEngine) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reg indicates an expected call of Reg.
func (mr *MockRestEndpointRegisterRecorder) Reg(arg0 RestEngine) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reg", reflect.TypeOf((*MockRestEndpointRegister)(nil).Reg), arg0)
}
