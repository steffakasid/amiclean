// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	context "context"

	cloudtrail "github.com/aws/aws-sdk-go-v2/service/cloudtrail"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// CloudTrailClient is an autogenerated mock type for the CloudTrailClient type
type CloudTrailClient struct {
	mock.Mock
}

type CloudTrailClient_Expecter struct {
	mock *mock.Mock
}

func (_m *CloudTrailClient) EXPECT() *CloudTrailClient_Expecter {
	return &CloudTrailClient_Expecter{mock: &_m.Mock}
}

// LookupEvents provides a mock function with given fields: ctx, params, optFns
func (_m *CloudTrailClient) LookupEvents(ctx context.Context, params *cloudtrail.LookupEventsInput, optFns ...func(*cloudtrail.Options)) (*cloudtrail.LookupEventsOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *cloudtrail.LookupEventsOutput
	if rf, ok := ret.Get(0).(func(context.Context, *cloudtrail.LookupEventsInput, ...func(*cloudtrail.Options)) *cloudtrail.LookupEventsOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cloudtrail.LookupEventsOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *cloudtrail.LookupEventsInput, ...func(*cloudtrail.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CloudTrailClient_LookupEvents_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LookupEvents'
type CloudTrailClient_LookupEvents_Call struct {
	*mock.Call
}

// LookupEvents is a helper method to define mock.On call
//  - ctx context.Context
//  - params *cloudtrail.LookupEventsInput
//  - optFns ...func(*cloudtrail.Options)
func (_e *CloudTrailClient_Expecter) LookupEvents(ctx interface{}, params interface{}, optFns ...interface{}) *CloudTrailClient_LookupEvents_Call {
	return &CloudTrailClient_LookupEvents_Call{Call: _e.mock.On("LookupEvents",
		append([]interface{}{ctx, params}, optFns...)...)}
}

func (_c *CloudTrailClient_LookupEvents_Call) Run(run func(ctx context.Context, params *cloudtrail.LookupEventsInput, optFns ...func(*cloudtrail.Options))) *CloudTrailClient_LookupEvents_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]func(*cloudtrail.Options), len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(func(*cloudtrail.Options))
			}
		}
		run(args[0].(context.Context), args[1].(*cloudtrail.LookupEventsInput), variadicArgs...)
	})
	return _c
}

func (_c *CloudTrailClient_LookupEvents_Call) Return(_a0 *cloudtrail.LookupEventsOutput, _a1 error) *CloudTrailClient_LookupEvents_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// NewCloudTrailClient creates a new instance of CloudTrailClient. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewCloudTrailClient(t testing.TB) *CloudTrailClient {
	mock := &CloudTrailClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}