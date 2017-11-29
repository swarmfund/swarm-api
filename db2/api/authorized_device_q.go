package api

import (
	"database/sql"
	"fmt"

	sq "github.com/lann/squirrel"
)

const tableAuthorizedDevice = "authorized_device"

var authorizedDeviceSelect = sq.Select("*").From(tableAuthorizedDevice)
var authorizedDeviceInsert = sq.Insert(tableAuthorizedDevice)
var authorizedDeviceUpdate = sq.Update(tableAuthorizedDevice)

type AuthorizedDeviceQI interface {
	Create(device *AuthorizedDevice) error
	ByWalletID(dest []*AuthorizedDevice, walletId int64) (err error)
	ByFingerprint(fingerprint string) (result *AuthorizedDevice, err error)
	UpdateLastLoginTime(device *AuthorizedDevice) (err error)
}

type AuthorizedDeviceQ struct {
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) AuthorizedDevice() AuthorizedDeviceQI {
	return &AuthorizedDeviceQ{
		parent: q,
		sql:    authorizedDeviceSelect,
	}
}

func (q *AuthorizedDeviceQ) Create(device *AuthorizedDevice) error {
	sql := authorizedDeviceInsert.SetMap(map[string]interface{}{
		"wallet_id":   device.WalletID,
		"fingerprint": device.Fingerprint,
		"details":     device.Details,
	})

	_, err := q.parent.Exec(sql)
	return err
}

func (q *AuthorizedDeviceQ) ByWalletID(dest []*AuthorizedDevice, walletId int64) (err error) {
	q.sql = q.sql.Where("wallet_id = ?", walletId)
	err = q.parent.Select(&dest, q.sql)

	return err
}

func (q *AuthorizedDeviceQ) ByFingerprint(fingerprint string) (*AuthorizedDevice, error) {
	q.sql = q.sql.Where("fingerprint = ?", fingerprint)
	var result AuthorizedDevice
	err := q.parent.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (q *AuthorizedDeviceQ) UpdateLastLoginTime(device *AuthorizedDevice) (err error) {
	sql := fmt.Sprintf(`update %s
		set last_login_at = timestamp 'now'
		where id = $1`, tableAuthorizedDevice)

	_, err = q.parent.DB.Exec(sql, device.Id)
	return err
}
