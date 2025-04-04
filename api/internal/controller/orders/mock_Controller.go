// Code generated by mockery v2.45.1. DO NOT EDIT.

package orders

import (
	context "context"
	model "omg/api/internal/model"

	mock "github.com/stretchr/testify/mock"
)

// MockController is an autogenerated mock type for the Controller type
type MockController struct {
	mock.Mock
}

// CreateOrder provides a mock function with given fields: _a0, _a1
func (_m *MockController) CreateOrder(_a0 context.Context, _a1 model.CreateOrderInput) (model.Order, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrder")
	}

	var r0 model.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.CreateOrderInput) (model.Order, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.CreateOrderInput) model.Order); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(model.Order)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.CreateOrderInput) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOrderStatus provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockController) UpdateOrderStatus(_a0 context.Context, _a1 int64, _a2 model.OrderStatus) (model.Order, error) {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrderStatus")
	}

	var r0 model.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, model.OrderStatus) (model.Order, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, model.OrderStatus) model.Order); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(model.Order)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, model.OrderStatus) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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
