package api

import (
	"database/sql"

	"encoding/json"

	"strings"

	sq "github.com/lann/squirrel"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/api/tfa"
)

var walletSelect = sq.Select(
	"w.*",
	"et.confirmed as verified",
	"r.address as recovery_address",
	"r.wallet_id as recovery_wallet_id",
	"r.salt as recovery_salt").
	From("wallets w").
	// TODO make just join
	LeftJoin("recoveries r on w.email = r.wallet").
	Join("email_tokens et on w.wallet_id = et.wallet_id")

var walletInsert = sq.Insert("wallets")
var walletUpdate = sq.Update("wallets")

const (
	tableWallets                    = "wallets"
	tableRecoveries                 = "recoveries"
	tableWalletsLimit               = 10
	walletsKDFFkeyConstraint        = `wallets_kdf_fkey`
	walletsWalletIDKeyConstraint    = `wallets_wallet_id_key`
	recoveriesWalletIDKeyConstraint = `recoveries_wallet_id_unique_constraint`
)

var (
	ErrWalletsKDFViolated      = errors.New("wallets_kdf_fkey violated")
	ErrWalletsWalletIDViolated = errors.New("wallets_wallet_id_key violated")
	ErrWalletsConflict         = errors.New("wallet already exists")
	ErrRecoveriesConflict      = errors.New("recovery already exists")
)

type RecoveryKeychain struct {
	Email    string        `db:"wallet"`
	Salt     string        `db:"salt"`
	Keychain string        `db:"keychain"`
	Address  types.Address `db:"address"`
	WalletID string        `db:"wallet_id"`
}

//go:generate mockery -case underscore -name WalletQI
type WalletQI interface {
	New() WalletQI
	Transaction(func(WalletQI) error) error

	// KDF
	KDFByVersion(int64) (*data.KDF, error)
	CreateWalletKDF(data.WalletKDF) error
	KDFByEmail(string) (*data.KDF, error)
	UpdateWalletKDF(data.WalletKDF) error

	// TODO belongs to TFAQI, it's here for transaction implementation reasons
	// CreatePasswordFactor will mutate factor with ID
	CreatePasswordFactor(walletID string, factor *tfa.Password) error
	// DeletePasswordFactor assumes there is single password factor per wallet
	DeletePasswordFactor(walletID string) error
	// TODO also belongs somewhere else, here for same reasons
	CreateRecovery(RecoveryKeychain) error

	// Create expected to set wallet.ID on successful create
	// May throw:
	// * ErrWalletsKDFViolated if KDF version is invalid
	// * ErrWalletsWalletIDViolated if wallet_id is not unique
	Create(wallet *Wallet) error

	// LoadWallet
	ByEmail(username string) (*Wallet, error)
	ByWalletID(walletId string) (*Wallet, error)
	DeleteWallets(walletIDs []string) error
	ByWalletOrRecoveryID(walletId string) (*Wallet, bool, error)

	ByCurrentAccountID(accountID string) (*Wallet, error)
	ByAccountID(address types.Address) (*Wallet, error)

	Update(w *Wallet) error

	// DEPRECATED
	Page(uint64) WalletQI
	// DEPRECATED
	ByState(uint64) WalletQI
	// DEPRECATED
	Select() ([]Wallet, error)
}

type WalletQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *WalletQ) New() WalletQI {
	return &WalletQ{
		parent: &Q{
			q.parent.Clone(),
		},
		sql: walletSelect,
	}
}

func (q *WalletQ) KDFByVersion(version int64) (*data.KDF, error) {
	var result data.KDF
	stmt := sq.Select("*").From("kdf").Where("version = ?", version)
	err := q.parent.Get(&result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (q *WalletQ) KDFByEmail(email string) (*data.KDF, error) {
	var result data.KDF
	stmt := sq.
		Select("kdf.version", "kdf.algorithm", "kdf.bits",
			"kdf.n", "kdf.r", "kdf.p", "kdf_wallets.salt").
		From("kdf_wallets").
		Join("kdf on kdf.version = kdf_wallets.version").
		Where("kdf_wallets.wallet = ?", email)
	err := q.parent.Get(&result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (q *WalletQ) CreateWalletKDF(kdf data.WalletKDF) error {
	stmt := sq.Insert("kdf_wallets").SetMap(map[string]interface{}{
		"salt":    kdf.Salt,
		"version": kdf.Version,
		"wallet":  kdf.Wallet,
	})
	_, err := q.parent.Exec(stmt)
	return err
}

func (q *WalletQ) UpdateWalletKDF(kdf data.WalletKDF) error {
	stmt := sq.Update("kdf_wallets").SetMap(map[string]interface{}{
		"salt":    kdf.Salt,
		"version": kdf.Version,
		"wallet":  kdf.Wallet,
	})
	_, err := q.parent.Exec(stmt)
	return err
}
func (q *WalletQ) Transaction(fn func(q WalletQI) error) (err error) {
	return q.parent.Transaction(func() error {
		return fn(q)
	})
}

func (q *WalletQ) CreateRecovery(recovery RecoveryKeychain) error {
	stmt := sq.Insert(tableRecoveries).
		SetMap(map[string]interface{}{
			"wallet":        recovery.Email,
			"salt":          recovery.Salt,
			"keychain_data": recovery.Keychain,
			"wallet_id":     recovery.WalletID,
			"address":       recovery.Address,
		})
	_, err := q.parent.Exec(stmt)
	if err != nil {
		pqerr, ok := errors.Cause(err).(*pq.Error)
		if ok {
			if pqerr.Constraint == recoveriesWalletIDKeyConstraint {
				return ErrRecoveriesConflict
			}
		}
	}
	return err
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

func (q *WalletQ) ByWalletID(walletID string) (*Wallet, error) {
	var result Wallet
	stmt := walletSelect.Where("w.wallet_id = ?", walletID)

	err := q.parent.Get(&result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *WalletQ) ByWalletOrRecoveryID(walletID string) (*Wallet, bool, error) {
	var result Wallet
	stmt := walletSelect.Where("w.wallet_id = ? OR r.wallet_id = ?", walletID, walletID)

	err := q.parent.Get(&result, stmt)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}

	return &result, result.RecoveryWalletID == walletID, err
}

func (q *WalletQ) DeleteWallets(walletIDs []string) error {
	if len(walletIDs) == 0 {
		return nil
	}

	sqq := sq.Delete(tableWallets).Where(sq.Eq{"wallet_id": walletIDs})
	_, err := q.parent.Exec(sqq)
	return err
}

func (q *WalletQ) ByState(state uint64) WalletQI {
	if q.Err != nil {
		return q
	}
	conditions := []string{}

	if state&uint64(types.WalletStateNotVerified) != 0 {
		conditions = append(conditions, "et.confirmed = false")
	}

	if state&uint64(types.WalletStateVerified) != 0 {
		conditions = append(conditions, "et.confirmed = true")
	}

	if len(conditions) > 0 {
		q.sql = q.sql.Where(strings.Join(conditions, " OR "))
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

func (q *WalletQ) Update(w *Wallet) error {
	stmt := walletUpdate.SetMap(map[string]interface{}{
		"wallet_id":          w.WalletId,
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
