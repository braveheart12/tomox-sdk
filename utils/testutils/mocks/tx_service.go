// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
import mock "github.com/stretchr/testify/mock"

import types "github.com/tomochain/tomoxsdk/types"

// TxService is an autogenerated mock type for the TxService type
type TxService struct {
	mock.Mock
}

// GetCustomTxSendOptions provides a mock function with given fields: w
func (_m *TxService) GetCustomTxSendOptions(w *types.Wallet) *bind.TransactOpts {
	ret := _m.Called(w)

	var r0 *bind.TransactOpts
	if rf, ok := ret.Get(0).(func(*types.Wallet) *bind.TransactOpts); ok {
		r0 = rf(w)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bind.TransactOpts)
		}
	}

	return r0
}

// GetTxCallOptions provides a mock function with given fields:
func (_m *TxService) GetTxCallOptions() *bind.CallOpts {
	ret := _m.Called()

	var r0 *bind.CallOpts
	if rf, ok := ret.Get(0).(func() *bind.CallOpts); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bind.CallOpts)
		}
	}

	return r0
}

// GetTxDefaultSendOptions provides a mock function with given fields:
func (_m *TxService) GetTxDefaultSendOptions() (*bind.TransactOpts, error) {
	ret := _m.Called()

	var r0 *bind.TransactOpts
	if rf, ok := ret.Get(0).(func() *bind.TransactOpts); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bind.TransactOpts)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTxSendOptions provides a mock function with given fields:
func (_m *TxService) GetTxSendOptions() (*bind.TransactOpts, error) {
	ret := _m.Called()

	var r0 *bind.TransactOpts
	if rf, ok := ret.Get(0).(func() *bind.TransactOpts); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bind.TransactOpts)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetTxSender provides a mock function with given fields: w
func (_m *TxService) SetTxSender(w *types.Wallet) {
	_m.Called(w)
}
