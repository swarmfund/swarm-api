// Code generated by mockery v1.0.0
package mocks

import data "gitlab.com/swarmfund/api/internal/data"
import mock "github.com/stretchr/testify/mock"

// EmailTokensQ is an autogenerated mock type for the EmailTokensQ type
type EmailTokensQ struct {
	mock.Mock
}

// Create provides a mock function with given fields: walletID, token
func (_m *EmailTokensQ) Create(walletID string, token string) error {
	ret := _m.Called(walletID, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(walletID, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: walletID
func (_m *EmailTokensQ) Get(walletID string) (*data.EmailToken, error) {
	ret := _m.Called(walletID)

	var r0 *data.EmailToken
	if rf, ok := ret.Get(0).(func(string) *data.EmailToken); ok {
		r0 = rf(walletID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*data.EmailToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(walletID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUnsent provides a mock function with given fields:
func (_m *EmailTokensQ) GetUnsent() ([]data.EmailToken, error) {
	ret := _m.Called()

	var r0 []data.EmailToken
	if rf, ok := ret.Get(0).(func() []data.EmailToken); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]data.EmailToken)
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

// MarkSent provides a mock function with given fields: tid
func (_m *EmailTokensQ) MarkSent(tid int64) error {
	ret := _m.Called(tid)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(tid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MarkUnsent provides a mock function with given fields: tid
func (_m *EmailTokensQ) MarkUnsent(tid int64) error {
	ret := _m.Called(tid)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(tid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// New provides a mock function with given fields:
func (_m *EmailTokensQ) New() data.EmailTokensQ {
	ret := _m.Called()

	var r0 data.EmailTokensQ
	if rf, ok := ret.Get(0).(func() data.EmailTokensQ); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(data.EmailTokensQ)
		}
	}

	return r0
}

// Verify provides a mock function with given fields: walletID, token
func (_m *EmailTokensQ) Verify(walletID string, token string) (bool, error) {
	ret := _m.Called(walletID, token)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(walletID, token)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(walletID, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
