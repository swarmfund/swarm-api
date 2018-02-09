// Code generated by mockery v1.0.0
package mocks

import api "gitlab.com/swarmfund/api/db2/api"
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

// ByWalletIDOrRecovery provides a mock function with given fields: walletId
func (_m *WalletQI) ByWalletIDOrRecovery(walletId string) (*api.Wallet, error) {
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

// RecoveryByWalletID provides a mock function with given fields: recoveryWalletID
func (_m *WalletQI) RecoveryByWalletID(recoveryWalletID string) (*api.RecoveryKeychain, error) {
	ret := _m.Called(recoveryWalletID)

	var r0 *api.RecoveryKeychain
	if rf, ok := ret.Get(0).(func(string) *api.RecoveryKeychain); ok {
		r0 = rf(recoveryWalletID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.RecoveryKeychain)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(recoveryWalletID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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
