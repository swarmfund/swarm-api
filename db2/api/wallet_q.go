package api

import (
	"database/sql"

	"encoding/json"

	sq "github.com/lann/squirrel"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/api/tfa"
)

var walletSelect = sq.Select(
	"w.*",
	"ow.organization_address",
	"ow.id is null as detached",
	"et.confirmed as verified").
	From("wallets w").
	LeftJoin("organization_wallets ow ON w.id = ow.wallet_id").
	Join("email_tokens et on w.wallet_id = et.wallet_id")

var walletInsert = sq.Insert("wallets")
var walletUpdate = sq.Update("wallets")

const (
	tableWalletsLimit            = 2
	walletsKDFFkeyConstraint     = `wallets_kdf_fkey`
	walletsWalletIDKeyConstraint = `wallets_wallet_id_key`
)

var (
	ErrWalletsKDFViolated      = errors.New("wallets_kdf_fkey violated")
	ErrWalletsWalletIDViolated = errors.New("wallets_wallet_id_key violated")
	ErrWalletsConflict         = errors.New("wallet already exists")
)

//go:generate mockery -case underscore -name WalletQI
type WalletQI interface {
	New() WalletQI
	Transaction(func(WalletQI) error) error

	// TODO belongs to TFAQI, it's here for transaction implementation reasons
	// CreatePasswordFactor will mutate factor with ID
	CreatePasswordFactor(walletID string, factor *tfa.Password) error
	// DeletePasswordFactor assumes there is single password factor per wallet
	DeletePasswordFactor(walletID string) error

	// Create expected to set wallet.ID on successful create
	// May throw:
	// * ErrWalletsKDFViolated if KDF version is invalid
	// * ErrWalletsWalletIDViolated if wallet_id is not unique
	Create(wallet *Wallet) error
	CreateOrganizationAttachment(wid int64) error
	UpdateOrganizationAttachment(wid int64, address, operation string) error
	OrganizationWatcherCursor() (string, error)

	// LoadWallet
	ByEmail(username string) (*Wallet, error)
	ByWalletID(walletId string) (*Wallet, error)
	// it's all S3 fault
	RecoveryWallet(lowercaseWalletID, username string) (*Wallet, error)

	ByID(id int64) (*Wallet, error)
	ByCurrentAccountID(accountID string) (*Wallet, error)
	ByAccountID(address types.Address) (*Wallet, error)

	Verify(id int64) error

	Update(w *Wallet) error

	Delete(id int64) error
	// delete all wallets except provided one
	SetActive(accountID, walletID string) error

	Page(uint64) WalletQI
	ByState(uint64) WalletQI
	Select() ([]Wallet, error)
}

type WalletQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *WalletQ) New() WalletQI {
	return &WalletQ{
		parent: q.parent,
		sql:    walletSelect,
	}
}

func (q *WalletQ) Transaction(fn func(q WalletQI) error) (err error) {
	return q.parent.Transaction(func() error {
		return fn(q)
	})
}

func (q *WalletQ) CreatePasswordFactor(walletID string, factor *tfa.Password) error {
	model := Backend{
		BackendType: types.WalletFactorPassword,
		WalletID:    walletID,
		Priority:    1,
	}
	details, err := json.Marshal(&factor)
	if err != nil {
		return errors.Wrap(err, "failed to marshal details")
	}
	model.Details = details

	stmt := sq.Insert(tfaBackendTable).
		SetMap(map[string]interface{}{
			"wallet_id": model.WalletID,
			"details":   model.Details,
			"priority":  model.Priority,
			"backend":   model.BackendType,
		}).Suffix("returning id")

	var factorID int64
	err = q.parent.Get(&factorID, stmt)
	if err != nil {
		cause := errors.Cause(err)
		pqerr, ok := cause.(*pq.Error)
		if ok {
			if pqerr.Constraint == tfaWalletFactorPasswordConstraint {
				return ErrWalletBackendConflict
			}
		}
	}
	mutatedFactor, err := tfa.NewPasswordFromDB(factorID, details)
	if err != nil {
		return errors.Wrap(err, "failed to mutate factor")
	}

	*factor = *mutatedFactor

	return err
}

func (q *WalletQ) DeletePasswordFactor(walletID string) error {
	stmt := sq.Delete(tfaBackendTable).
		Where("wallet_id = ?", walletID).
		Where("backend = ?", types.WalletFactorPassword)

	_, err := q.parent.Exec(stmt)
	return err
}

func (q *WalletQ) Page(page uint64) WalletQI {
	if q.Err != nil {
		return q
	}
	q.sql = q.sql.Offset(tableWalletsLimit * (page - 1)).Limit(tableWalletsLimit)
	return q
}

func (q *WalletQ) Create(w *Wallet) error {
	stmt := walletInsert.SetMap(map[string]interface{}{
		"wallet_id":          w.WalletId,
		"account_id":         w.AccountID,
		"current_account_id": w.CurrentAccountID,
		"email":              w.Username,
		"salt":               w.Salt,
		"kdf_id":             w.KDF,
		"keychain_data":      w.KeychainData,
	}).Suffix("returning id")

	err := q.parent.Get(&(w.Id), stmt)
	if err != nil {
		cause := errors.Cause(err)
		pqerr, ok := cause.(*pq.Error)
		if ok {
			if pqerr.Constraint == walletsKDFFkeyConstraint {
				return ErrWalletsKDFViolated
			}
			if pqerr.Constraint == walletsWalletIDKeyConstraint {
				return ErrWalletsWalletIDViolated
			}
		}
	}
	return err
}

func (q *WalletQ) CreateOrganizationAttachment(wid int64) error {
	stmt := sq.Insert("organization_wallets").
		SetMap(map[string]interface{}{
			"wallet_id": wid,
		})
	_, err := q.parent.Exec(stmt)
	return err
}

func (q *WalletQ) UpdateOrganizationAttachment(wid int64, address, operation string) error {
	stmt := sq.Update("organization_wallets").SetMap(map[string]interface{}{
		"organization_address": address,
		"operation":            operation,
	}).Where("wallet_id = ?", wid)
	_, err := q.parent.Exec(stmt)
	return err
}

func (q *WalletQ) OrganizationWatcherCursor() (string, error) {
	result := ""
	stmt := sq.
		Select("operation").
		From("organization_wallets").
		OrderBy("operation desc").
		Limit(1)
	err := q.parent.Get(&result, stmt)
	if err == sql.ErrNoRows {
		return "0", nil
	}
	return result, err
}

func (q *WalletQ) ByEmail(email string) (*Wallet, error) {
	var result Wallet
	stmt := walletSelect.Where("w.email = ?", email)

	err := q.parent.Get(&result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *WalletQ) ByCurrentAccountID(accountID string) (*Wallet, error) {
	var result Wallet
	stmt := walletSelect.Where("w.current_account_id = ?", accountID)
	err := q.parent.Get(&result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (q *WalletQ) ByAccountID(address types.Address) (*Wallet, error) {
	result := &Wallet{}
	stmt := walletSelect.Where("w.account_id = ?", address)
	err := q.parent.Get(result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

func (q *WalletQ) ByID(id int64) (*Wallet, error) {
	result := &Wallet{}
	stmt := walletSelect.Where("w.id = ?", id)
	err := q.parent.Get(result, stmt)
	return result, err
}

func (q *WalletQ) RecoveryWallet(lowercaseWalletID, username string) (*Wallet, error) {
	result := &Wallet{}
	stmt := walletSelect.Where("lower(w.wallet_id) = lower(?) and w.username = ?", lowercaseWalletID, username)
	err := q.parent.Get(result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return result, err
}

func (q *WalletQ) ByWalletID(walletId string) (*Wallet, error) {
	var result Wallet
	stmt := walletSelect.Where("w.wallet_id = ?", walletId)

	err := q.parent.Get(&result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *WalletQ) ByState(state uint64) WalletQI {
	if q.Err != nil {
		return q
	}
	// TODO proper state filters
	if state&uint64(1) != 0 {
		q.sql = q.sql.Where("et.confirmed= ?", false)
	} else if state&uint64(2) != 0 {
		q.sql = q.sql.Where("et.confirmed = ?", true)
	}

	return q
}

func (q *WalletQ) Select() ([]Wallet, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	var result []Wallet
	q.Err = q.parent.Select(&result, q.sql)
	return result, q.Err
}

func (q *WalletQ) SetActive(accountID string, walletID string) error {
	stmt := sq.Delete("wallets").
		Where("account_id = ?", accountID).
		Where("wallet_id != ?", walletID)

	_, err := q.parent.Exec(stmt)
	return err
}

func (q *WalletQ) Delete(id int64) error {
	if q.Err != nil {
		return q.Err
	}

	sqq := sq.Delete("wallets").Where("id = ?", id)
	_, q.Err = q.parent.Exec(sqq)

	return q.Err
}

func (q *WalletQ) Verify(id int64) error {
	if q.Err != nil {
		return q.Err
	}

	sqq := walletUpdate.Set("verified", true).Where("id = ?", id)
	_, q.Err = q.parent.Exec(sqq)

	return q.Err
}

func (q *WalletQ) Update(w *Wallet) error {
	stmt := walletUpdate.SetMap(map[string]interface{}{
		"wallet_id":          w.WalletId,
		"salt":               w.Salt,
		"current_account_id": w.CurrentAccountID,
		"keychain_data":      w.KeychainData,
	}).Where("id = ?", w.Id)

	_, err := q.parent.Exec(stmt)
	if err != nil {
		cause := errors.Cause(err)
		pqerr, ok := cause.(*pq.Error)
		if ok {
			if pqerr.Constraint == walletsKDFFkeyConstraint {
				return ErrWalletsKDFViolated
			}
			if pqerr.Constraint == walletsWalletIDKeyConstraint {
				return ErrWalletsWalletIDViolated
			}
		}
	}

	return err
}
