// Code generated by mockery v2.45.1. DO NOT EDIT.

package system

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockController is an autogenerated mock type for the Controller type
type MockController struct {
	mock.Mock
}

// CheckReadiness provides a mock function with given fields: ctx
func (_m *MockController) CheckReadiness(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for CheckReadiness")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockController creates a new instance of MockController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockController(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockController {
	mock := &MockController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
