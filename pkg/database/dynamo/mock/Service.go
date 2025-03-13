// Code generated by mockery v2.52.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// DeleteItem provides a mock function with given fields: ctx, key
func (_m *Service) DeleteItem(ctx context.Context, key map[string]interface{}) error {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for DeleteItem")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetItem provides a mock function with given fields: ctx, key
func (_m *Service) GetItem(ctx context.Context, key map[string]interface{}) (map[string]interface{}, error) {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for GetItem")
	}

	var r0 map[string]interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}) (map[string]interface{}, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}) map[string]interface{}); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, map[string]interface{}) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutItem provides a mock function with given fields: ctx, item
func (_m *Service) PutItem(ctx context.Context, item map[string]interface{}) error {
	ret := _m.Called(ctx, item)

	if len(ret) == 0 {
		panic("no return value specified for PutItem")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateItem provides a mock function with given fields: ctx, key, update
func (_m *Service) UpdateItem(ctx context.Context, key map[string]interface{}, update map[string]interface{}) error {
	ret := _m.Called(ctx, key, update)

	if len(ret) == 0 {
		panic("no return value specified for UpdateItem")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, map[string]interface{}) error); ok {
		r0 = rf(ctx, key, update)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
