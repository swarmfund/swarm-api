// Code generated by mockery v1.0.0
package mocks

import api "gitlab.com/swarmfund/api/db2/api"
import data "gitlab.com/swarmfund/api/internal/data"
import mock "github.com/stretchr/testify/mock"
import tfa "gitlab.com/swarmfund/api/tfa"
import types "gitlab.com/swarmfund/api/internal/types"

// WalletQI is an autogenerated mock type for the WalletQI type
type WalletQI struct {
	mock.Mock
}

// ByAccountID provides a mock function with given fields: address
func (_m *WalletQI) ByAccountID(address types.Address) (*api.Wallet, error) {
	ret := _m.Called(address)

	var r0 *api.Wallet
	if rf, ok := ret.Get(0).(func(types.Address) *api.Wallet); ok {
		r0 = rf(address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Wallet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.Address) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ByCurrentAccountID provides a mock function with given fields: accountID
func (_m *WalletQI) ByCurrentAccountID(accountID string) (*api.Wallet, error) {
	ret := _m.Called(accountID)

	var r0 *api.Wallet
	if rf, ok := ret.Get(0).(func(string) *api.Wallet); ok {
		r0 = rf(accountID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Wallet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(accountID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ByEmail provides a mock function with given fields: username
func (_m *WalletQI) ByEmail(username string) (*api.Wallet, error) {
	ret := _m.Called(username)

	var r0 *api.Wallet
	if rf, ok := ret.Get(0).(func(string) *api.Wallet); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Wallet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ByState provides a mock function with given fields: _a0
func (_m *WalletQI) ByState(_a0 uint64) api.WalletQI {
	ret := _m.Called(_a0)

	var r0 api.WalletQI
	if rf, ok := ret.Get(0).(func(uint64) api.WalletQI); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.WalletQI)
		}
	}

	return r0
}

// ByWalletID provides a mock function with given fields: walletId
func (_m *WalletQI) ByWalletID(walletId string) (*api.Wallet, error) {
	ret := _m.Called(walletId)

	var r0 *api.Wallet
	if rf, ok := ret.Get(0).(func(string) *api.Wallet); ok {
		r0 = rf(walletId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Wallet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(walletId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ByWalletOrRecoveryID provides a mock function with given fields: walletId
func (_m *WalletQI) ByWalletOrRecoveryID(walletId string) (*api.Wallet, bool, error) {
	ret := _m.Called(walletId)

	var r0 *api.Wallet
	if rf, ok := ret.Get(0).(func(string) *api.Wallet); ok {
		r0 = rf(walletId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Wallet)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(walletId)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(walletId)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Create provides a mock function with given fields: wallet
func (_m *WalletQI) Create(wallet *api.Wallet) error {
	ret := _m.Called(wallet)

	var r0 error
	if rf, ok := ret.Get(0).(func(*api.Wallet) error); ok {
		r0 = rf(wallet)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreatePasswordFactor provides a mock function with given fields: walletID, factor
func (_m *WalletQI) CreatePasswordFactor(walletID string, factor *tfa.Password) error {
	ret := _m.Called(walletID, factor)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *tfa.Password) error); ok {
		r0 = rf(walletID, factor)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateRecovery provides a mock function with given fields: _a0
func (_m *WalletQI) CreateRecovery(_a0 api.RecoveryKeychain) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(api.RecoveryKeychain) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateReferral provides a mock function with given fields: referrer, referral
func (_m *WalletQI) CreateReferral(referrer types.Address, referral types.Address) error {
	ret := _m.Called(referrer, referral)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Address, types.Address) error); ok {
		r0 = rf(referrer, referral)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateWalletKDF provides a mock function with given fields: _a0
func (_m *WalletQI) CreateWalletKDF(_a0 data.WalletKDF) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(data.WalletKDF) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: walletID
func (_m *WalletQI) Delete(walletID string) error {
	ret := _m.Called(walletID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(walletID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePasswordFactor provides a mock function with given fields: walletID
func (_m *WalletQI) DeletePasswordFactor(walletID string) error {
	ret := _m.Called(walletID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(walletID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteWallets provides a mock function with given fields: walletIDs
func (_m *WalletQI) DeleteWallets(walletIDs []string) error {
	ret := _m.Called(walletIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(walletIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// KDFByEmail provides a mock function with given fields: _a0
func (_m *WalletQI) KDFByEmail(_a0 string) (*data.KDF, error) {
	ret := _m.Called(_a0)

	var r0 *data.KDF
	if rf, ok := ret.Get(0).(func(string) *data.KDF); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*data.KDF)
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

// KDFByVersion provides a mock function with given fields: _a0
func (_m *WalletQI) KDFByVersion(_a0 int64) (*data.KDF, error) {
	ret := _m.Called(_a0)

	var r0 *data.KDF
	if rf, ok := ret.Get(0).(func(int64) *data.KDF); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*data.KDF)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// New provides a mock function with given fields:
func (_m *WalletQI) New() api.WalletQI {
	ret := _m.Called()

	var r0 api.WalletQI
	if rf, ok := ret.Get(0).(func() api.WalletQI); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.WalletQI)
		}
	}

	return r0
}

// Page provides a mock function with given fields: _a0
func (_m *WalletQI) Page(_a0 uint64) api.WalletQI {
	ret := _m.Called(_a0)

	var r0 api.WalletQI
	if rf, ok := ret.Get(0).(func(uint64) api.WalletQI); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.WalletQI)
		}
	}

	return r0
}

// Select provides a mock function with given fields:
func (_m *WalletQI) Select() ([]api.Wallet, error) {
	ret := _m.Called()

	var r0 []api.Wallet
	if rf, ok := ret.Get(0).(func() []api.Wallet); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]api.Wallet)
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

// Transaction provides a mock function with given fields: _a0
func (_m *WalletQI) Transaction(_a0 func(api.WalletQI) error) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(func(api.WalletQI) error) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: w
func (_m *WalletQI) Update(w *api.Wallet) error {
	ret := _m.Called(w)

	var r0 error
	if rf, ok := ret.Get(0).(func(*api.Wallet) error); ok {
		r0 = rf(w)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateWalletKDF provides a mock function with given fields: _a0
func (_m *WalletQI) UpdateWalletKDF(_a0 data.WalletKDF) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(data.WalletKDF) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
