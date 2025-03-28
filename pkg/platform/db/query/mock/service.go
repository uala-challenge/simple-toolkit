// Code generated by mockery v2.52.3. DO NOT EDIT.

package mocks

import (
	context "context"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	mock "github.com/stretchr/testify/mock"

	types "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Apply provides a mock function with given fields: ctx, item
func (_m *Service) Apply(ctx context.Context, item *dynamodb.QueryInput) ([]map[string]types.AttributeValue, error) {
	ret := _m.Called(ctx, item)

	if len(ret) == 0 {
		panic("no return value specified for Apply")
	}

	var r0 []map[string]types.AttributeValue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dynamodb.QueryInput) ([]map[string]types.AttributeValue, error)); ok {
		return rf(ctx, item)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dynamodb.QueryInput) []map[string]types.AttributeValue); ok {
		r0 = rf(ctx, item)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]map[string]types.AttributeValue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dynamodb.QueryInput) error); ok {
		r1 = rf(ctx, item)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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
