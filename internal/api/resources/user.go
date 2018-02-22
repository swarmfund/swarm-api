package resources

import (
	"time"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/types"
)

type (
	User struct {
		Type       types.UserType `json:"type"`
		ID         types.Address  `json:"id"`
		Attributes UserAttributes `json:"attributes"`
	}
	UserAttributes struct {
		Email           string             `json:"email"`
		State           types.UserState    `json:"state"`
		KYCSequence     int64              `json:"kyc_sequence"`
		RejectReason    string             `json:"reject_reason"`
		RecoveryAddress types.Address      `json:"recovery_address"`
		CreatedAt       time.Time          `json:"created_at"`
		AirdropState    types.AirdropState `json:"airdrop_state"`
	}
)

// NewUser populates user resource from db record
func NewUser(user *api.User) User {
	// worth place for this, but there is nothing better
	if user.AirdropState == nil {
		state := types.AirdropStateEligible
		// checking if user is eligible
		if !user.IsAirdropEligible() {
			state = types.AirdropStateNotEligible
		}
		user.AirdropState = &state
	}
	return User{
		Type: user.UserType,
		ID:   user.Address,
		Attributes: UserAttributes{
			Email:           user.Email,
			State:           user.State,
			KYCSequence:     user.KYCSequence,
			RejectReason:    user.RejectReason,
			RecoveryAddress: user.RecoveryAddress,
			AirdropState:    *user.AirdropState,
			CreatedAt:       user.CreatedAt,
		},
	}
}
