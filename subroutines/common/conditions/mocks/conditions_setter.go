// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConditionsSetter is an autogenerated mock type for the ConditionsSetter type
type ConditionsSetter struct {
	mock.Mock
}

type ConditionsSetter_Expecter struct {
	mock *mock.Mock
}

func (_m *ConditionsSetter) EXPECT() *ConditionsSetter_Expecter {
	return &ConditionsSetter_Expecter{mock: &_m.Mock}
}

// SetFalse provides a mock function with given fields: objectMeta, _a1, reason, message
func (_m *ConditionsSetter) SetFalse(objectMeta v1.ObjectMeta, _a1 *[]v1.Condition, reason string, message string) {
	_m.Called(objectMeta, _a1, reason, message)
}

// ConditionsSetter_SetFalse_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetFalse'
type ConditionsSetter_SetFalse_Call struct {
	*mock.Call
}

// SetFalse is a helper method to define mock.On call
//   - objectMeta v1.ObjectMeta
//   - _a1 *[]v1.Condition
//   - reason string
//   - message string
func (_e *ConditionsSetter_Expecter) SetFalse(objectMeta interface{}, _a1 interface{}, reason interface{}, message interface{}) *ConditionsSetter_SetFalse_Call {
	return &ConditionsSetter_SetFalse_Call{Call: _e.mock.On("SetFalse", objectMeta, _a1, reason, message)}
}

func (_c *ConditionsSetter_SetFalse_Call) Run(run func(objectMeta v1.ObjectMeta, _a1 *[]v1.Condition, reason string, message string)) *ConditionsSetter_SetFalse_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(v1.ObjectMeta), args[1].(*[]v1.Condition), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *ConditionsSetter_SetFalse_Call) Return() *ConditionsSetter_SetFalse_Call {
	_c.Call.Return()
	return _c
}

// SetTrue provides a mock function with given fields: objectMeta, _a1, reason, message
func (_m *ConditionsSetter) SetTrue(objectMeta v1.ObjectMeta, _a1 *[]v1.Condition, reason string, message string) {
	_m.Called(objectMeta, _a1, reason, message)
}

// ConditionsSetter_SetTrue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetTrue'
type ConditionsSetter_SetTrue_Call struct {
	*mock.Call
}

// SetTrue is a helper method to define mock.On call
//   - objectMeta v1.ObjectMeta
//   - _a1 *[]v1.Condition
//   - reason string
//   - message string
func (_e *ConditionsSetter_Expecter) SetTrue(objectMeta interface{}, _a1 interface{}, reason interface{}, message interface{}) *ConditionsSetter_SetTrue_Call {
	return &ConditionsSetter_SetTrue_Call{Call: _e.mock.On("SetTrue", objectMeta, _a1, reason, message)}
}

func (_c *ConditionsSetter_SetTrue_Call) Run(run func(objectMeta v1.ObjectMeta, _a1 *[]v1.Condition, reason string, message string)) *ConditionsSetter_SetTrue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(v1.ObjectMeta), args[1].(*[]v1.Condition), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *ConditionsSetter_SetTrue_Call) Return() *ConditionsSetter_SetTrue_Call {
	_c.Call.Return()
	return _c
}

type mockConstructorTestingTNewConditionsSetter interface {
	mock.TestingT
	Cleanup(func())
}

// NewConditionsSetter creates a new instance of ConditionsSetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConditionsSetter(t mockConstructorTestingTNewConditionsSetter) *ConditionsSetter {
	mock := &ConditionsSetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}