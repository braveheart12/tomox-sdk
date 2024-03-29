// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import bson "github.com/globalsign/mgo/bson"
import common "github.com/ethereum/go-ethereum/common"

import mock "github.com/stretchr/testify/mock"
import types "github.com/tomochain/tomoxsdk/types"

// PairDao is an autogenerated mock type for the PairDao type
type PairDao struct {
	mock.Mock
}

// Create provides a mock function with given fields: o
func (_m *PairDao) Create(o *types.Pair) error {
	ret := _m.Called(o)

	var r0 error
	if rf, ok := ret.Get(0).(func(*types.Pair) error); ok {
		r0 = rf(o)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetActivePairs provides a mock function with given fields:
func (_m *PairDao) GetActivePairs() ([]*types.Pair, error) {
	ret := _m.Called()

	var r0 []*types.Pair
	if rf, ok := ret.Get(0).(func() []*types.Pair); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*types.Pair)
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

// GetAll provides a mock function with given fields:
func (_m *PairDao) GetAll() ([]types.Pair, error) {
	ret := _m.Called()

	var r0 []types.Pair
	if rf, ok := ret.Get(0).(func() []types.Pair); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Pair)
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

// GetByID provides a mock function with given fields: id
func (_m *PairDao) GetByID(id bson.ObjectId) (*types.Pair, error) {
	ret := _m.Called(id)

	var r0 *types.Pair
	if rf, ok := ret.Get(0).(func(bson.ObjectId) *types.Pair); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Pair)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bson.ObjectId) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByName provides a mock function with given fields: name
func (_m *PairDao) GetByName(name string) (*types.Pair, error) {
	ret := _m.Called(name)

	var r0 *types.Pair
	if rf, ok := ret.Get(0).(func(string) *types.Pair); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Pair)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByTokenAddress provides a mock function with given fields: baseToken, quoteToken
func (_m *PairDao) GetByTokenAddress(baseToken common.Address, quoteToken common.Address) (*types.Pair, error) {
	ret := _m.Called(baseToken, quoteToken)

	var r0 *types.Pair
	if rf, ok := ret.Get(0).(func(common.Address, common.Address) *types.Pair); ok {
		r0 = rf(baseToken, quoteToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Pair)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(common.Address, common.Address) error); ok {
		r1 = rf(baseToken, quoteToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByTokenSymbols provides a mock function with given fields: baseTokenSymbol, quoteTokenSymbol
func (_m *PairDao) GetByTokenSymbols(baseTokenSymbol string, quoteTokenSymbol string) (*types.Pair, error) {
	ret := _m.Called(baseTokenSymbol, quoteTokenSymbol)

	var r0 *types.Pair
	if rf, ok := ret.Get(0).(func(string, string) *types.Pair); ok {
		r0 = rf(baseTokenSymbol, quoteTokenSymbol)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Pair)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(baseTokenSymbol, quoteTokenSymbol)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
