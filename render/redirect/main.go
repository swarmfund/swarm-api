package redirect

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type ActionType string

var (
	ActionRecoveryRequest      ActionType = "recovery_request"
	ActionSignup               ActionType = "signup"
	ActionSkrillDepositSuccess ActionType = "skrill_deposit_success"
	ActionSkrillDepositCancel  ActionType = "skrill_deposit_cancel"
	ActionPendingWallet        ActionType = "pending_wallet"

	// generic redirects

	SignupPayload = Payload{
		Status: http.StatusOK,
		Action: ActionSignup,
	}

	ServerError = Payload{
		Status: http.StatusInternalServerError,
	}

	NotFound = Payload{
		Status: http.StatusNotFound,
	}

	Unavailable = Payload{
		Status: http.StatusServiceUnavailable,
	}

	// recovery request

	RecoveryRequestAlreadyUploaded = Payload{
		Status: http.StatusBadGateway,
		Action: ActionRecoveryRequest,
		Data: map[string]interface{}{
			"reason": "already uploaded",
		},
	}

	RecoveryRequestShowCode = func(username, code string, isFulfilled bool) *Payload {
		return &Payload{
			Status: http.StatusOK,
			Action: ActionRecoveryRequest,
			Data: map[string]interface{}{
				"username":     username,
				"code":         code,
				"is_fulfilled": isFulfilled,
			},
		}
	}

	// skrill deposit

	SkrillDepositSuccess = Payload{
		Status: http.StatusOK,
		Action: ActionSkrillDepositSuccess,
	}

	SkrillDepositCancel = Payload{
		Status: http.StatusOK,
		Action: ActionSkrillDepositCancel,
	}

	// 3rd party KYC

	PendingWallet = func(accountID string) *Payload {
		return &Payload{
			Status: http.StatusOK,
			Action: ActionPendingWallet,
			Data: map[string]interface{}{
				"account_id": accountID,
			},
		}
	}
)

type Payload struct {
	Status int                    `json:"status"`
	Action ActionType             `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

func (p *Payload) Encode() (string, error) {
	bytes, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(bytes)
	if err != nil {
		return "", err
	}
	return encoded, nil
}
