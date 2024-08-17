// Code generated by mockery v2.44.2. DO NOT EDIT.

package common

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// MockScheduler is an autogenerated mock type for the Scheduler type
type MockScheduler struct {
	mock.Mock
}

type MockScheduler_Expecter struct {
	mock *mock.Mock
}

func (_m *MockScheduler) EXPECT() *MockScheduler_Expecter {
	return &MockScheduler_Expecter{mock: &_m.Mock}
}

// CreateOrUpdateOneTimeTask provides a mock function with given fields: tag, at, task
func (_m *MockScheduler) CreateOrUpdateOneTimeTask(tag string, at time.Time, task func() error) error {
	ret := _m.Called(tag, at, task)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrUpdateOneTimeTask")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, time.Time, func() error) error); ok {
		r0 = rf(tag, at, task)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockScheduler_CreateOrUpdateOneTimeTask_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateOrUpdateOneTimeTask'
type MockScheduler_CreateOrUpdateOneTimeTask_Call struct {
	*mock.Call
}

// CreateOrUpdateOneTimeTask is a helper method to define mock.On call
//   - tag string
//   - at time.Time
//   - task func() error
func (_e *MockScheduler_Expecter) CreateOrUpdateOneTimeTask(tag interface{}, at interface{}, task interface{}) *MockScheduler_CreateOrUpdateOneTimeTask_Call {
	return &MockScheduler_CreateOrUpdateOneTimeTask_Call{Call: _e.mock.On("CreateOrUpdateOneTimeTask", tag, at, task)}
}

func (_c *MockScheduler_CreateOrUpdateOneTimeTask_Call) Run(run func(tag string, at time.Time, task func() error)) *MockScheduler_CreateOrUpdateOneTimeTask_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(time.Time), args[2].(func() error))
	})
	return _c
}

func (_c *MockScheduler_CreateOrUpdateOneTimeTask_Call) Return(_a0 error) *MockScheduler_CreateOrUpdateOneTimeTask_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockScheduler_CreateOrUpdateOneTimeTask_Call) RunAndReturn(run func(string, time.Time, func() error) error) *MockScheduler_CreateOrUpdateOneTimeTask_Call {
	_c.Call.Return(run)
	return _c
}

// CreateOrUpdateRecurringTask provides a mock function with given fields: tag, cron, task
func (_m *MockScheduler) CreateOrUpdateRecurringTask(tag string, cron string, task func() error) error {
	ret := _m.Called(tag, cron, task)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrUpdateRecurringTask")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, func() error) error); ok {
		r0 = rf(tag, cron, task)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockScheduler_CreateOrUpdateRecurringTask_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateOrUpdateRecurringTask'
type MockScheduler_CreateOrUpdateRecurringTask_Call struct {
	*mock.Call
}

// CreateOrUpdateRecurringTask is a helper method to define mock.On call
//   - tag string
//   - cron string
//   - task func() error
func (_e *MockScheduler_Expecter) CreateOrUpdateRecurringTask(tag interface{}, cron interface{}, task interface{}) *MockScheduler_CreateOrUpdateRecurringTask_Call {
	return &MockScheduler_CreateOrUpdateRecurringTask_Call{Call: _e.mock.On("CreateOrUpdateRecurringTask", tag, cron, task)}
}

func (_c *MockScheduler_CreateOrUpdateRecurringTask_Call) Run(run func(tag string, cron string, task func() error)) *MockScheduler_CreateOrUpdateRecurringTask_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(func() error))
	})
	return _c
}

func (_c *MockScheduler_CreateOrUpdateRecurringTask_Call) Return(_a0 error) *MockScheduler_CreateOrUpdateRecurringTask_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockScheduler_CreateOrUpdateRecurringTask_Call) RunAndReturn(run func(string, string, func() error) error) *MockScheduler_CreateOrUpdateRecurringTask_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteTask provides a mock function with given fields: tag
func (_m *MockScheduler) DeleteTask(tag string) error {
	ret := _m.Called(tag)

	if len(ret) == 0 {
		panic("no return value specified for DeleteTask")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(tag)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockScheduler_DeleteTask_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteTask'
type MockScheduler_DeleteTask_Call struct {
	*mock.Call
}

// DeleteTask is a helper method to define mock.On call
//   - tag string
func (_e *MockScheduler_Expecter) DeleteTask(tag interface{}) *MockScheduler_DeleteTask_Call {
	return &MockScheduler_DeleteTask_Call{Call: _e.mock.On("DeleteTask", tag)}
}

func (_c *MockScheduler_DeleteTask_Call) Run(run func(tag string)) *MockScheduler_DeleteTask_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockScheduler_DeleteTask_Call) Return(_a0 error) *MockScheduler_DeleteTask_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockScheduler_DeleteTask_Call) RunAndReturn(run func(string) error) *MockScheduler_DeleteTask_Call {
	_c.Call.Return(run)
	return _c
}

// GetTaskNextRun provides a mock function with given fields: tag
func (_m *MockScheduler) GetTaskNextRun(tag string) string {
	ret := _m.Called(tag)

	if len(ret) == 0 {
		panic("no return value specified for GetTaskNextRun")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(tag)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockScheduler_GetTaskNextRun_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTaskNextRun'
type MockScheduler_GetTaskNextRun_Call struct {
	*mock.Call
}

// GetTaskNextRun is a helper method to define mock.On call
//   - tag string
func (_e *MockScheduler_Expecter) GetTaskNextRun(tag interface{}) *MockScheduler_GetTaskNextRun_Call {
	return &MockScheduler_GetTaskNextRun_Call{Call: _e.mock.On("GetTaskNextRun", tag)}
}

func (_c *MockScheduler_GetTaskNextRun_Call) Run(run func(tag string)) *MockScheduler_GetTaskNextRun_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockScheduler_GetTaskNextRun_Call) Return(_a0 string) *MockScheduler_GetTaskNextRun_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockScheduler_GetTaskNextRun_Call) RunAndReturn(run func(string) string) *MockScheduler_GetTaskNextRun_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockScheduler creates a new instance of MockScheduler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockScheduler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockScheduler {
	mock := &MockScheduler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
