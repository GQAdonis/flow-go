// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// ApprovalValidator is an autogenerated mock type for the ApprovalValidator type
type ApprovalValidator struct {
	mock.Mock
}

// Validate provides a mock function with given fields: approval
func (_m *ApprovalValidator) Validate(approval *flow.ResultApproval) error {
	ret := _m.Called(approval)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.ResultApproval) error); ok {
		r0 = rf(approval)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}