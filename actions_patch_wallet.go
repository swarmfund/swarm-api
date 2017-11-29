package api

import (
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
	"gitlab.com/swarmfund/api/render/problem"
)

type PatchWalletRequest struct {
	// Detached flag denoting that wallet is used in organization multi-sign flow,
	// could not be updated if wallet is already attached to organization
	Detached *bool `json:"detached"`
}

type PatchWalletAction struct {
	Action

	WalletID int64
	Request  PatchWalletRequest

	Wallet *api.Wallet
}

func (action *PatchWalletAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.loadWallet,
		action.checkAllowed,
		action.updateWallet,
		func() {
			hal.Render(action.W, problem.Success)
		},
	)
}

func (action *PatchWalletAction) loadParams() {
	action.WalletID = action.GetInt64("id")
	action.UnmarshalBody(&action.Request)
}

func (action *PatchWalletAction) loadWallet() {
	wallet, err := action.APIQ().Wallet().ByID(action.WalletID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get wallet")
		action.Err = &problem.ServerError
		return
	}

	if wallet == nil {
		action.Err = &problem.NotFound
	}

	if !wallet.Verified {
		// forbidden is probably a better response,
		// but just to not leak wallet state until we checked signature
		action.Err = &problem.NotFound
	}

	action.Wallet = wallet
}

func (action *PatchWalletAction) checkAllowed() {
	action.checkSignerConstraints(
		SignedBy(action.Wallet.CurrentAccountID),
	)
}

func (action *PatchWalletAction) updateWallet() {
	q := action.APIQ()
	err := q.GetRepo().Begin()
	if err != nil {
		action.Log.WithError(err).Error("failed to begin tx")
		action.Err = &problem.ServerError
		return
	}

	// check if we are actually updating `detached`
	if detached := action.Request.Detached; detached != nil && action.Wallet.Detached != *detached {
		if *detached == true {
			// don't try to optimize it and check for current multi-sign state,
			// in case of ingest race condition wallet enters non-deterministic state
			action.SetInvalidField("detached", errors.New("could not be updated"))
			return
		}

		// check if user exists
		user, err := q.Users().ByAddress(action.Wallet.AccountID)
		if err != nil {
			action.Log.WithError(err).Error("failed to get user")
			action.Err = &problem.ServerError
			return
		}
		if user != nil {
			action.SetInvalidField("detached", errors.New("could not be updated"))
			return
		}
		err = q.Wallet().CreateOrganizationAttachment(int64(action.Wallet.Id))
		if err != nil {
			action.Log.WithError(err).Error("failed to create org attachment")
			action.Err = &problem.ServerError
			return
		}
	}

	err = q.GetRepo().Commit()
	if err != nil {
		action.Log.WithError(err).Error("failed to commit tx")
		action.Err = &problem.ServerError
		return
	}
}
