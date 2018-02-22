package resources

import (
	"time"

	"regexp"

	"github.com/pkg/errors"
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

func isEmailAirdropEligible(email string) bool {
	pattern := `(?i).+@(163.com|qq.com|126.com|189.cn|139.com|sina.com|aliyun.com|xinjiyuan99.com)`
	blacklisted, err := regexp.MatchString(pattern, email)
	if err != nil {
		panic(errors.Wrap(err, "blacklist check failed"))
	}
	return !blacklisted
}

// NewUser populates user resource from db record
func NewUser(user *api.User) User {
	// worth place for this, but there is nothing better
	if user.AirdropState == nil {
		state := types.AirdropStateEligible
		// checking if email domain is blacklisted
		if !isEmailAirdropEligible(user.Email) {
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
