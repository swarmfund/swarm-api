// Code generated by mockery v1.0.0
package mocks

import mock "github.com/stretchr/testify/mock"
import types "gitlab.com/swarmfund/api/internal/types"

// Blobs is an autogenerated mock type for the Blobs type
type Blobs struct {
	mock.Mock
}

// Create provides a mock function with given fields: address, blob
func (_m *Blobs) Create(address types.Address, blob *types.Blob) error {
	ret := _m.Called(address, blob)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Address, *types.Blob) error); ok {
		r0 = rf(address, blob)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: id
func (_m *Blobs) Get(id string) (*types.Blob, error) {
	ret := _m.Called(id)

	var r0 *types.Blob
	if rf, ok := ret.Get(0).(func(string) *types.Blob); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Blob)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
