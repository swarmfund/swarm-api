package api

import (
	"strconv"
	"time"

	"github.com/guregu/null"
	"gitlab.com/swarmfund/api/internal/lorem"
)

type RecoveryRequest struct {
	ID               int64     `db:"id" json:"id"`
	WalletID         int64     `db:"wallet_id" json:"wallet_id"` // actual db wallet id not `wallet_id`
	Username         string    `db:"username" json:"username"`
	AccountID        string    `db:"account_id" json:"account_id"`
	EmailToken       string    `db:"email_token" json:"-"`
	Code             string    `db:"code" json:"code"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	SentAt           null.Time `db:"sent_at" json:"sent_at"`
	CodeShownAt      null.Time `db:"code_shown_at" json:"code_shown_at"`
	UploadedAt       null.Time `db:"uploaded_at" json:"uploaded_at"`
	RecoveryWalletID *string   `db:"recovery_wallet_id" json:"recovery_wallet_id"`
}

func NewRecoveryRequest(wallet *Wallet) (*RecoveryRequest, error) {
	accountID := wallet.AccountID
	if wallet.OrganizationAddress != nil {
		accountID = wallet.AccountID
	}
	return &RecoveryRequest{
		WalletID:   wallet.Id,
		AccountID:  accountID,
		Username:   wallet.Username,
		EmailToken: lorem.Token(),
		Code:       lorem.Token(),
	}, nil
}

func (r RecoveryRequest) PagingToken() string {
	return strconv.FormatInt(r.ID, 10)
}

func (r *RecoveryRequest) CodeShown() bool {
	return r.CodeShownAt.Valid
}

func (r *RecoveryRequest) IsUploaded() bool {
	return r.UploadedAt.Valid
}
