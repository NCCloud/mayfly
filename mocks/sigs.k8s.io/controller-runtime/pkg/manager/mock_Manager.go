// Code generated by mockery v2.52.2. DO NOT EDIT.

package manager

import (
	cache "sigs.k8s.io/controller-runtime/pkg/cache"
	client "sigs.k8s.io/controller-runtime/pkg/client"

	config "sigs.k8s.io/controller-runtime/pkg/config"

	context "context"

	healthz "sigs.k8s.io/controller-runtime/pkg/healthz"

	http "net/http"

	logr "github.com/go-logr/logr"

	manager "sigs.k8s.io/controller-runtime/pkg/manager"

	meta "k8s.io/apimachinery/pkg/api/meta"

	mock "github.com/stretchr/testify/mock"

	record "k8s.io/client-go/tools/record"

	rest "k8s.io/client-go/rest"

	runtime "k8s.io/apimachinery/pkg/runtime"

	webhook "sigs.k8s.io/controller-runtime/pkg/webhook"
)

// MockManager is an autogenerated mock type for the Manager type
type MockManager struct {
	mock.Mock
}

type MockManager_Expecter struct {
	mock *mock.Mock
}

func (_m *MockManager) EXPECT() *MockManager_Expecter {
	return &MockManager_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: _a0
func (_m *MockManager) Add(_a0 manager.Runnable) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(manager.Runnable) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockManager_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type MockManager_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - _a0 manager.Runnable
func (_e *MockManager_Expecter) Add(_a0 interface{}) *MockManager_Add_Call {
	return &MockManager_Add_Call{Call: _e.mock.On("Add", _a0)}
}

func (_c *MockManager_Add_Call) Run(run func(_a0 manager.Runnable)) *MockManager_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(manager.Runnable))
	})
	return _c
}

func (_c *MockManager_Add_Call) Return(_a0 error) *MockManager_Add_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_Add_Call) RunAndReturn(run func(manager.Runnable) error) *MockManager_Add_Call {
	_c.Call.Return(run)
	return _c
}

// AddHealthzCheck provides a mock function with given fields: name, check
func (_m *MockManager) AddHealthzCheck(name string, check healthz.Checker) error {
	ret := _m.Called(name, check)

	if len(ret) == 0 {
		panic("no return value specified for AddHealthzCheck")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, healthz.Checker) error); ok {
		r0 = rf(name, check)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockManager_AddHealthzCheck_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddHealthzCheck'
type MockManager_AddHealthzCheck_Call struct {
	*mock.Call
}

// AddHealthzCheck is a helper method to define mock.On call
//   - name string
//   - check healthz.Checker
func (_e *MockManager_Expecter) AddHealthzCheck(name interface{}, check interface{}) *MockManager_AddHealthzCheck_Call {
	return &MockManager_AddHealthzCheck_Call{Call: _e.mock.On("AddHealthzCheck", name, check)}
}

func (_c *MockManager_AddHealthzCheck_Call) Run(run func(name string, check healthz.Checker)) *MockManager_AddHealthzCheck_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(healthz.Checker))
	})
	return _c
}

func (_c *MockManager_AddHealthzCheck_Call) Return(_a0 error) *MockManager_AddHealthzCheck_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_AddHealthzCheck_Call) RunAndReturn(run func(string, healthz.Checker) error) *MockManager_AddHealthzCheck_Call {
	_c.Call.Return(run)
	return _c
}

// AddMetricsServerExtraHandler provides a mock function with given fields: path, handler
func (_m *MockManager) AddMetricsServerExtraHandler(path string, handler http.Handler) error {
	ret := _m.Called(path, handler)

	if len(ret) == 0 {
		panic("no return value specified for AddMetricsServerExtraHandler")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, http.Handler) error); ok {
		r0 = rf(path, handler)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockManager_AddMetricsServerExtraHandler_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddMetricsServerExtraHandler'
type MockManager_AddMetricsServerExtraHandler_Call struct {
	*mock.Call
}

// AddMetricsServerExtraHandler is a helper method to define mock.On call
//   - path string
//   - handler http.Handler
func (_e *MockManager_Expecter) AddMetricsServerExtraHandler(path interface{}, handler interface{}) *MockManager_AddMetricsServerExtraHandler_Call {
	return &MockManager_AddMetricsServerExtraHandler_Call{Call: _e.mock.On("AddMetricsServerExtraHandler", path, handler)}
}

func (_c *MockManager_AddMetricsServerExtraHandler_Call) Run(run func(path string, handler http.Handler)) *MockManager_AddMetricsServerExtraHandler_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(http.Handler))
	})
	return _c
}

func (_c *MockManager_AddMetricsServerExtraHandler_Call) Return(_a0 error) *MockManager_AddMetricsServerExtraHandler_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_AddMetricsServerExtraHandler_Call) RunAndReturn(run func(string, http.Handler) error) *MockManager_AddMetricsServerExtraHandler_Call {
	_c.Call.Return(run)
	return _c
}

// AddReadyzCheck provides a mock function with given fields: name, check
func (_m *MockManager) AddReadyzCheck(name string, check healthz.Checker) error {
	ret := _m.Called(name, check)

	if len(ret) == 0 {
		panic("no return value specified for AddReadyzCheck")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, healthz.Checker) error); ok {
		r0 = rf(name, check)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockManager_AddReadyzCheck_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddReadyzCheck'
type MockManager_AddReadyzCheck_Call struct {
	*mock.Call
}

// AddReadyzCheck is a helper method to define mock.On call
//   - name string
//   - check healthz.Checker
func (_e *MockManager_Expecter) AddReadyzCheck(name interface{}, check interface{}) *MockManager_AddReadyzCheck_Call {
	return &MockManager_AddReadyzCheck_Call{Call: _e.mock.On("AddReadyzCheck", name, check)}
}

func (_c *MockManager_AddReadyzCheck_Call) Run(run func(name string, check healthz.Checker)) *MockManager_AddReadyzCheck_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(healthz.Checker))
	})
	return _c
}

func (_c *MockManager_AddReadyzCheck_Call) Return(_a0 error) *MockManager_AddReadyzCheck_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_AddReadyzCheck_Call) RunAndReturn(run func(string, healthz.Checker) error) *MockManager_AddReadyzCheck_Call {
	_c.Call.Return(run)
	return _c
}

// Elected provides a mock function with no fields
func (_m *MockManager) Elected() <-chan struct{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Elected")
	}

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// MockManager_Elected_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Elected'
type MockManager_Elected_Call struct {
	*mock.Call
}

// Elected is a helper method to define mock.On call
func (_e *MockManager_Expecter) Elected() *MockManager_Elected_Call {
	return &MockManager_Elected_Call{Call: _e.mock.On("Elected")}
}

func (_c *MockManager_Elected_Call) Run(run func()) *MockManager_Elected_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockManager_Elected_Call) Return(_a0 <-chan struct{}) *MockManager_Elected_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_Elected_Call) RunAndReturn(run func() <-chan struct{}) *MockManager_Elected_Call {
	_c.Call.Return(run)
	return _c
}

// GetAPIReader provides a mock function with no fields
func (_m *MockManager) GetAPIReader() client.Reader {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAPIReader")
	}

	var r0 client.Reader
	if rf, ok := ret.Get(0).(func() client.Reader); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(client.Reader)
		}
	}

	return r0
}

// MockManager_GetAPIReader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAPIReader'
type MockManager_GetAPIReader_Call struct {
	*mock.Call
}

// GetAPIReader is a helper method to define mock.On call
func (_e *MockManager_Expecter) GetAPIReader() *MockManager_GetAPIReader_Call {
	return &MockManager_GetAPIReader_Call{Call: _e.mock.On("GetAPIReader")}
}

func (_c *MockManager_GetAPIReader_Call) Run(run func()) *MockManager_GetAPIReader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockManager_GetAPIReader_Call) Return(_a0 client.Reader) *MockManager_GetAPIReader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_GetAPIReader_Call) RunAndReturn(run func() client.Reader) *MockManager_GetAPIReader_Call {
	_c.Call.Return(run)
	return _c
}

// GetCache provides a mock function with no fields
func (_m *MockManager) GetCache() cache.Cache {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetCache")
	}

	var r0 cache.Cache
	if rf, ok := ret.Get(0).(func() cache.Cache); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cache.Cache)
		}
	}

	return r0
}

// MockManager_GetCache_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCache'
type MockManager_GetCache_Call struct {
	*mock.Call
}

// GetCache is a helper method to define mock.On call
func (_e *MockManager_Expecter) GetCache() *MockManager_GetCache_Call {
	return &MockManager_GetCache_Call{Call: _e.mock.On("GetCache")}
}

func (_c *MockManager_GetCache_Call) Run(run func()) *MockManager_GetCache_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockManager_GetCache_Call) Return(_a0 cache.Cache) *MockManager_GetCache_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_GetCache_Call) RunAndReturn(run func() cache.Cache) *MockManager_GetCache_Call {
	_c.Call.Return(run)
	return _c
}

// GetClient provides a mock function with no fields
func (_m *MockManager) GetClient() client.Client {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetClient")
	}

	var r0 client.Client
	if rf, ok := ret.Get(0).(func() client.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(client.Client)
		}
	}

	return r0
}

// MockManager_GetClient_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetClient'
type MockManager_GetClient_Call struct {
	*mock.Call
}

// GetClient is a helper method to define mock.On call
func (_e *MockManager_Expecter) GetClient() *MockManager_GetClient_Call {
	return &MockManager_GetClient_Call{Call: _e.mock.On("GetClient")}
}

func (_c *MockManager_GetClient_Call) Run(run func()) *MockManager_GetClient_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockManager_GetClient_Call) Return(_a0 client.Client) *MockManager_GetClient_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_GetClient_Call) RunAndReturn(run func() client.Client) *MockManager_GetClient_Call {
	_c.Call.Return(run)
	return _c
}

// GetConfig provides a mock function with no fields
func (_m *MockManager) GetConfig() *rest.Config {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetConfig")
	}

	var r0 *rest.Config
	if rf, ok := ret.Get(0).(func() *rest.Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rest.Config)
		}
	}

	return r0
}

// MockManager_GetConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConfig'
type MockManager_GetConfig_Call struct {
	*mock.Call
}

// GetConfig is a helper method to define mock.On call
func (_e *MockManager_Expecter) GetConfig() *MockManager_GetConfig_Call {
	return &MockManager_GetConfig_Call{Call: _e.mock.On("GetConfig")}
}

func (_c *MockManager_GetConfig_Call) Run(run func()) *MockManager_GetConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockManager_GetConfig_Call) Return(_a0 *rest.Config) *MockManager_GetConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_GetConfig_Call) RunAndReturn(run func() *rest.Config) *MockManager_GetConfig_Call {
	_c.Call.Return(run)
	return _c
}

// GetControllerOptions provides a mock function with no fields
func (_m *MockManager) GetControllerOptions() config.Controller {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetControllerOptions")
	}

	var r0 config.Controller
	if rf, ok := ret.Get(0).(func() config.Controller); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(config.Controller)
	}

	return r0
}

// MockManager_GetControllerOptions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetControllerOptions'
type MockManager_GetControllerOptions_Call struct {
	*mock.Call
}

// GetControllerOptions is a helper method to define mock.On call
func (_e *MockManager_Expecter) GetControllerOptions() *MockManager_GetControllerOptions_Call {
	return &MockManager_GetControllerOptions_Call{Call: _e.mock.On("GetControllerOptions")}
}

func (_c *MockManager_GetControllerOptions_Call) Run(run func()) *MockManager_GetControllerOptions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockManager_GetControllerOptions_Call) Return(_a0 config.Controller) *MockManager_GetControllerOptions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_GetControllerOptions_Call) RunAndReturn(run func() config.Controller) *MockManager_GetControllerOptions_Call {
	_c.Call.Return(run)
	return _c
}

// GetEventRecorderFor provides a mock function with given fields: name
func (_m *MockManager) GetEventRecorderFor(name string) record.EventRecorder {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for GetEventRecorderFor")
	}

	var r0 record.EventRecorder
	if rf, ok := ret.Get(0).(func(string) record.EventRecorder); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(record.EventRecorder)
		}
	}

	return r0
}

// MockManager_GetEventRecorderFor_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetEventRecorderFor'
type MockManager_GetEventRecorderFor_Call struct {
	*mock.Call
}

// GetEventRecorderFor is a helper method to define mock.On call
//   - name string
func (_e *MockManager_Expecter) GetEventRecorderFor(name interface{}) *MockManager_GetEventRecorderFor_Call {
	return &MockManager_GetEventRecorderFor_Call{Call: _e.mock.On("GetEventRecorderFor", name)}
}

func (_c *MockManager_GetEventRecorderFor_Call) Run(run func(name string)) *MockManager_GetEventRecorderFor_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockManager_GetEventRecorderFor_Call) Return(_a0 record.EventRecorder) *MockManager_GetEventRecorderFor_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_GetEventRecorderFor_Call) RunAndReturn(run func(string) record.EventRecorder) *MockManager_GetEventRecorderFor_Call {
	_c.Call.Return(run)
	return _c
}

// GetFieldIndexer provides a mock function with no fields
func (_m *MockManager) GetFieldIndexer() client.FieldIndexer {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetFieldIndexer")
	}

	var r0 client.FieldIndexer
	if rf, ok := ret.Get(0).(func() client.FieldIndexer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(client.FieldIndexer)
		}
	}

	return r0
}

// MockManager_GetFieldIndexer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFieldIndexer'
type MockManager_GetFieldIndexer_Call struct {
	*mock.Call
}

// GetFieldIndexer is a helper method to define mock.On call
func (_e *MockManager_Expecter) GetFieldIndexer() *MockManager_GetFieldIndexer_Call {
	return &MockManager_GetFieldIndexer_Call{Call: _e.mock.On("GetFieldIndexer")}
}

func (_c *MockManager_GetFieldIndexer_Call) Run(run func()) *MockManager_GetFieldIndexer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockManager_GetFieldIndexer_Call) Return(_a0 client.FieldIndexer) *MockManager_GetFieldIndexer_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_GetFieldIndexer_Call) RunAndReturn(run func() client.FieldIndexer) *MockManager_GetFieldIndexer_Call {
	_c.Call.Return(run)
	return _c
}

// GetHTTPClient provides a mock function with no fields
func (_m *MockManager) GetHTTPClient() *http.Client {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetHTTPClient")
	}

	var r0 *http.Client
	if rf, ok := ret.Get(0).(func() *http.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Client)
		}
	}

	return r0
}

// MockManager_GetHTTPClient_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetHTTPClient'
type MockManager_GetHTTPClient_Call struct {
	*mock.Call
}

// GetHTTPClient is a helper method to define mock.On call
func (_e *MockManager_Expecter) GetHTTPClient() *MockManager_GetHTTPClient_Call {
	return &MockManager_GetHTTPClient_Call{Call: _e.mock.On("GetHTTPClient")}
}

func (_c *MockManager_GetHTTPClient_Call) Run(run func()) *MockManager_GetHTTPClient_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockManager_GetHTTPClient_Call) Return(_a0 *http.Client) *MockManager_GetHTTPClient_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_GetHTTPClient_Call) RunAndReturn(run func() *http.Client) *MockManager_GetHTTPClient_Call {
	_c.Call.Return(run)
	return _c
}

// GetLogger provides a mock function with no fields
func (_m *MockManager) GetLogger() logr.Logger {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetLogger")
	}

	var r0 logr.Logger
	if rf, ok := ret.Get(0).(func() logr.Logger); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(logr.Logger)
	}

	return r0
}

// MockManager_GetLogger_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLogger'
type MockManager_GetLogger_Call struct {
	*mock.Call
}

// GetLogger is a helper method to define mock.On call
func (_e *MockManager_Expecter) GetLogger() *MockManager_GetLogger_Call {
	return &MockManager_GetLogger_Call{Call: _e.mock.On("GetLogger")}
}

func (_c *MockManager_GetLogger_Call) Run(run func()) *MockManager_GetLogger_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockManager_GetLogger_Call) Return(_a0 logr.Logger) *MockManager_GetLogger_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_GetLogger_Call) RunAndReturn(run func() logr.Logger) *MockManager_GetLogger_Call {
	_c.Call.Return(run)
	return _c
}

// GetRESTMapper provides a mock function with no fields
func (_m *MockManager) GetRESTMapper() meta.RESTMapper {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRESTMapper")
	}

	var r0 meta.RESTMapper
	if rf, ok := ret.Get(0).(func() meta.RESTMapper); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(meta.RESTMapper)
		}
	}

	return r0
}

// MockManager_GetRESTMapper_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRESTMapper'
type MockManager_GetRESTMapper_Call struct {
	*mock.Call
}

// GetRESTMapper is a helper method to define mock.On call
func (_e *MockManager_Expecter) GetRESTMapper() *MockManager_GetRESTMapper_Call {
	return &MockManager_GetRESTMapper_Call{Call: _e.mock.On("GetRESTMapper")}
}

func (_c *MockManager_GetRESTMapper_Call) Run(run func()) *MockManager_GetRESTMapper_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockManager_GetRESTMapper_Call) Return(_a0 meta.RESTMapper) *MockManager_GetRESTMapper_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_GetRESTMapper_Call) RunAndReturn(run func() meta.RESTMapper) *MockManager_GetRESTMapper_Call {
	_c.Call.Return(run)
	return _c
}

// GetScheme provides a mock function with no fields
func (_m *MockManager) GetScheme() *runtime.Scheme {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetScheme")
	}

	var r0 *runtime.Scheme
	if rf, ok := ret.Get(0).(func() *runtime.Scheme); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*runtime.Scheme)
		}
	}

	return r0
}

// MockManager_GetScheme_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetScheme'
type MockManager_GetScheme_Call struct {
	*mock.Call
}

// GetScheme is a helper method to define mock.On call
func (_e *MockManager_Expecter) GetScheme() *MockManager_GetScheme_Call {
	return &MockManager_GetScheme_Call{Call: _e.mock.On("GetScheme")}
}

func (_c *MockManager_GetScheme_Call) Run(run func()) *MockManager_GetScheme_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockManager_GetScheme_Call) Return(_a0 *runtime.Scheme) *MockManager_GetScheme_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_GetScheme_Call) RunAndReturn(run func() *runtime.Scheme) *MockManager_GetScheme_Call {
	_c.Call.Return(run)
	return _c
}

// GetWebhookServer provides a mock function with no fields
func (_m *MockManager) GetWebhookServer() webhook.Server {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetWebhookServer")
	}

	var r0 webhook.Server
	if rf, ok := ret.Get(0).(func() webhook.Server); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(webhook.Server)
		}
	}

	return r0
}

// MockManager_GetWebhookServer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWebhookServer'
type MockManager_GetWebhookServer_Call struct {
	*mock.Call
}

// GetWebhookServer is a helper method to define mock.On call
func (_e *MockManager_Expecter) GetWebhookServer() *MockManager_GetWebhookServer_Call {
	return &MockManager_GetWebhookServer_Call{Call: _e.mock.On("GetWebhookServer")}
}

func (_c *MockManager_GetWebhookServer_Call) Run(run func()) *MockManager_GetWebhookServer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockManager_GetWebhookServer_Call) Return(_a0 webhook.Server) *MockManager_GetWebhookServer_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_GetWebhookServer_Call) RunAndReturn(run func() webhook.Server) *MockManager_GetWebhookServer_Call {
	_c.Call.Return(run)
	return _c
}

// Start provides a mock function with given fields: ctx
func (_m *MockManager) Start(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockManager_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type MockManager_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockManager_Expecter) Start(ctx interface{}) *MockManager_Start_Call {
	return &MockManager_Start_Call{Call: _e.mock.On("Start", ctx)}
}

func (_c *MockManager_Start_Call) Run(run func(ctx context.Context)) *MockManager_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockManager_Start_Call) Return(_a0 error) *MockManager_Start_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockManager_Start_Call) RunAndReturn(run func(context.Context) error) *MockManager_Start_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockManager creates a new instance of MockManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockManager {
	mock := &MockManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
