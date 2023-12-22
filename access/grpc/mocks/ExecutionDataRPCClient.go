// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	executiondata "github.com/onflow/flow/protobuf/go/flow/executiondata"
	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// MockExecutionDataRPCClient is an autogenerated mock type for the ExecutionDataRPCClient type
type MockExecutionDataRPCClient struct {
	mock.Mock
}

// GetExecutionDataByBlockID provides a mock function with given fields: ctx, in, opts
func (_m *MockExecutionDataRPCClient) GetExecutionDataByBlockID(ctx context.Context, in *executiondata.GetExecutionDataByBlockIDRequest, opts ...grpc.CallOption) (*executiondata.GetExecutionDataByBlockIDResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *executiondata.GetExecutionDataByBlockIDResponse
	if rf, ok := ret.Get(0).(func(context.Context, *executiondata.GetExecutionDataByBlockIDRequest, ...grpc.CallOption) *executiondata.GetExecutionDataByBlockIDResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*executiondata.GetExecutionDataByBlockIDResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *executiondata.GetExecutionDataByBlockIDRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRegisterValues provides a mock function with given fields: ctx, in, opts
func (_m *MockExecutionDataRPCClient) GetRegisterValues(ctx context.Context, in *executiondata.GetRegisterValuesRequest, opts ...grpc.CallOption) (*executiondata.GetRegisterValuesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *executiondata.GetRegisterValuesResponse
	if rf, ok := ret.Get(0).(func(context.Context, *executiondata.GetRegisterValuesRequest, ...grpc.CallOption) *executiondata.GetRegisterValuesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*executiondata.GetRegisterValuesResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *executiondata.GetRegisterValuesRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscribeEvents provides a mock function with given fields: ctx, in, opts
func (_m *MockExecutionDataRPCClient) SubscribeEvents(ctx context.Context, in *executiondata.SubscribeEventsRequest, opts ...grpc.CallOption) (executiondata.ExecutionDataAPI_SubscribeEventsClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 executiondata.ExecutionDataAPI_SubscribeEventsClient
	if rf, ok := ret.Get(0).(func(context.Context, *executiondata.SubscribeEventsRequest, ...grpc.CallOption) executiondata.ExecutionDataAPI_SubscribeEventsClient); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(executiondata.ExecutionDataAPI_SubscribeEventsClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *executiondata.SubscribeEventsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubscribeExecutionData provides a mock function with given fields: ctx, in, opts
func (_m *MockExecutionDataRPCClient) SubscribeExecutionData(ctx context.Context, in *executiondata.SubscribeExecutionDataRequest, opts ...grpc.CallOption) (executiondata.ExecutionDataAPI_SubscribeExecutionDataClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 executiondata.ExecutionDataAPI_SubscribeExecutionDataClient
	if rf, ok := ret.Get(0).(func(context.Context, *executiondata.SubscribeExecutionDataRequest, ...grpc.CallOption) executiondata.ExecutionDataAPI_SubscribeExecutionDataClient); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(executiondata.ExecutionDataAPI_SubscribeExecutionDataClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *executiondata.SubscribeExecutionDataRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
