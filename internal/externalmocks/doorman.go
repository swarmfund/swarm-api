// Code generated by mockery v1.0.0
package externalmocks

import doorman "gitlab.com/tokend/go/doorman"
import http "net/http"
import mock "github.com/stretchr/testify/mock"
import resources "gitlab.com/tokend/go/resources"

// Doorman is an autogenerated mock type for the Doorman type
type Doorman struct {
	mock.Mock
}

// AccountSigners provides a mock function with given fields: _a0
func (_m *Doorman) AccountSigners(_a0 string) ([]resources.Signer, error) {
	ret := _m.Called(_a0)

	var r0 []resources.Signer
	if rf, ok := ret.Get(0).(func(string) []resources.Signer); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]resources.Signer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Check provides a mock function with given fields: _a0, _a1
func (_m *Doorman) Check(_a0 *http.Request, _a1 ...doorman.SignerConstraint) error {
	_va := make([]interface{}, len(_a1))
	for _i := range _a1 {
		_va[_i] = _a1[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*http.Request, ...doorman.SignerConstraint) error); ok {
		r0 = rf(_a0, _a1...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
