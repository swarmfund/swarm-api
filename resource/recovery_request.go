package resource

import (
	"time"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/horizon-connector"
	"github.com/guregu/null"
)

type RecoveryRequest struct {
	ID               int64      `json:"id"`
	Username         string     `json:"username"`
	AccountID        string     `json:"account_id"`
	Code             string     `json:"code"`
	CreatedAt        time.Time  `json:"created_at"`
	SentAt           null.Time  `json:"sent_at"`
	CodeShownAt      null.Time  `json:"code_shown_at"`
	UploadedAt       null.Time  `json:"uploaded_at"`
	RecoveryWalletID *string    `json:"recovery_wallet_id"`
	RecoverOp        *RecoverOp `json:"recover_op,omitempty"`
}

type RecoverOp struct {
	Account   string `json:"account"`
	OldSigner string `json:"old_signer"`
	NewSigner string `json:"new_signer"`
}

func (resource *RecoveryRequest) Populate(record *api.RecoveryRequest, op *horizon.RecoverOp) {
	resource.ID = record.ID
	resource.Username = record.Username
	resource.AccountID = record.AccountID
	resource.Code = record.Code
	resource.CreatedAt = record.CreatedAt
	resource.SentAt = record.SentAt
	resource.CodeShownAt = record.CodeShownAt
	resource.UploadedAt = record.UploadedAt
	resource.RecoveryWalletID = record.RecoveryWalletID
	if op != nil {
		resource.RecoverOp = &RecoverOp{
			Account:   op.AccountID,
			OldSigner: op.OldSigner,
			NewSigner: op.NewSigner,
		}
	}
}
