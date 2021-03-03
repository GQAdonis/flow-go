// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// ExecutionReceipts is an autogenerated mock type for the ExecutionReceipts type
type ExecutionReceipts struct {
	mock.Mock
}

// ByBlockID provides a mock function with given fields: blockID
func (_m *ExecutionReceipts) ByBlockID(blockID flow.Identifier) (*flow.ExecutionReceipt, error) {
	ret := _m.Called(blockID)

	var r0 *flow.ExecutionReceipt
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.ExecutionReceipt); ok {
		r0 = rf(blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.ExecutionReceipt)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ByBlockIDAllExecutionReceipts provides a mock function with given fields: blockID
func (_m *ExecutionReceipts) ByBlockIDAllExecutionReceipts(blockID flow.Identifier) ([]*flow.ExecutionReceipt, error) {
	ret := _m.Called(blockID)

	var r0 []*flow.ExecutionReceipt
	if rf, ok := ret.Get(0).(func(flow.Identifier) []*flow.ExecutionReceipt); ok {
		r0 = rf(blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*flow.ExecutionReceipt)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ByID provides a mock function with given fields: receiptID
func (_m *ExecutionReceipts) ByID(receiptID flow.Identifier) (*flow.ExecutionReceipt, error) {
	ret := _m.Called(receiptID)

	var r0 *flow.ExecutionReceipt
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.ExecutionReceipt); ok {
		r0 = rf(receiptID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.ExecutionReceipt)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(receiptID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Index provides a mock function with given fields: blockID, receiptID
func (_m *ExecutionReceipts) Index(blockID flow.Identifier, receiptID flow.Identifier) error {
	ret := _m.Called(blockID, receiptID)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier, flow.Identifier) error); ok {
		r0 = rf(blockID, receiptID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IndexByExecutor provides a mock function with given fields: receipt
func (_m *ExecutionReceipts) IndexByExecutor(receipt *flow.ExecutionReceipt) error {
	ret := _m.Called(receipt)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.ExecutionReceipt) error); ok {
		r0 = rf(receipt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Store provides a mock function with given fields: receipt
func (_m *ExecutionReceipts) Store(receipt *flow.ExecutionReceipt) error {
	ret := _m.Called(receipt)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.ExecutionReceipt) error); ok {
		r0 = rf(receipt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
