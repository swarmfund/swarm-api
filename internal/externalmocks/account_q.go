// Code generated by mockery v1.0.0. DO NOT EDIT.
package externalmocks

import mock "github.com/stretchr/testify/mock"
import resources "gitlab.com/tokend/go/resources"

// AccountQ is an autogenerated mock type for the AccountQ type
type AccountQ struct {
	mock.Mock
}

// Signers provides a mock function with given fields: address
func (_m *AccountQ) Signers(address string) ([]resources.Signer, error) {
	ret := _m.Called(address)

	var r0 []resources.Signer
	if rf, ok := ret.Get(0).(func(string) []resources.Signer); ok {
		r0 = rf(address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]resources.Signer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
