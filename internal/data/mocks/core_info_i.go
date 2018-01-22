// Code generated by mockery v1.0.0
package mocks

import mock "github.com/stretchr/testify/mock"

// CoreInfoI is an autogenerated mock type for the CoreInfoI type
type CoreInfoI struct {
	mock.Mock
}

// GetMasterAccountID provides a mock function with given fields:
func (_m *CoreInfoI) GetMasterAccountID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Passphrase provides a mock function with given fields:
func (_m *CoreInfoI) Passphrase() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// TXExpire provides a mock function with given fields:
func (_m *CoreInfoI) TXExpire() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}
