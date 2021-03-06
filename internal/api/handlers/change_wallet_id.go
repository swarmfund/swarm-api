package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/internal/secondfactor"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/api/tfa"
	"gitlab.com/tokend/go/doorman"
)

type (
	ChangeWalletIDRequest struct {
		resources.Wallet
		CurrentWalletID string `json:"-"`
	}
)

func NewChangeWalletIDRequest(r *http.Request) (ChangeWalletIDRequest, error) {
	request := ChangeWalletIDRequest{
		CurrentWalletID: chi.URLParam(r, "wallet-id"),
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r ChangeWalletIDRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.CurrentWalletID, Required),
		Field(&r.Wallet, Required),
	)
}

func ChangeWalletID(w http.ResponseWriter, r *http.Request) {
	request, err := NewChangeWalletIDRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// load wallet
	wallet, isRecovery, err := WalletQ(r).ByWalletOrRecoveryID(request.CurrentWalletID)
	if err != nil {
		Log(r).WithError(err).Error("failed to get wallet")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if wallet == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	// explicitly checking against recovery address for wallets that don't yet have account
	err = Doorman(r,
		doorman.SignatureOf(string(wallet.RecoveryAddress)),
		doorman.SignatureOf(string(wallet.CurrentAccountID)))
	if err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	// we are not forcing any 2fa checks if request is for recovery wallet
	if !isRecovery {
		// check if user knows password
		if err := secondfactor.NewConsumer(TFAQ(r)).WithTokenMixin("pwd-check").WithBackendType(types.WalletFactorPassword).Consume(r, wallet); err != nil {
			RenderFactorConsumeError(w, r, err)
			return
		}

		// ask for TOTP token if enabled
		if err := secondfactor.NewConsumer(TFAQ(r)).WithTokenMixin("totp-check").WithBackendType(types.WalletFactorTOTP).Consume(r, wallet); err != nil {
			RenderFactorConsumeError(w, r, err)
			return
		}
	}

	if isRecovery {
		// handling case when some one tries to recover non-existent account
		if err := performUserCreate(r, wallet); err != nil {
			switch errors.Cause(err) {
			case ErrWalletNotVerified:
				ape.RenderErr(w, movetoape.Forbidden("verification_required"))
				return
			case api.ErrUsersConflict:
				// it's ok, user already exists
			default:
				Log(r).WithError(err).Error("failed to create user")
				ape.RenderErr(w, problems.InternalError())
				return
			}
		}
	}

	factor := tfa.NewPasswordBackend(tfa.PasswordDetails{
		Salt:         request.Data.Relationships.Factor.Data.Attributes.Salt,
		AccountID:    types.Address(request.Data.Relationships.Factor.Data.Attributes.AccountID),
		KeychainData: request.Data.Relationships.Factor.Data.Attributes.KeychainData,
	})

	// update wallet
	wallet.WalletId = request.Data.ID
	wallet.KeychainData = request.Data.Attributes.KeychainData
	wallet.CurrentAccountID = types.Address(request.Data.Attributes.AccountID)
	// TODO transaction is not working. Error on horizon submition still makes commit!!!!!!!!!!
	err = WalletQ(r).Transaction(func(q api.WalletQI) error {
		// update wallet
		if err = q.Update(wallet); err != nil {
			return errors.Wrap(err, "failed to update wallet")
		}

		// update wallet kdf
		if err := q.UpdateWalletKDF(data.WalletKDF{
			Wallet:  wallet.Username,
			Version: request.Data.Relationships.KDF.Data.ID,
			Salt:    request.Data.Attributes.Salt,
		}); err != nil {
			return errors.Wrap(err, "failed to update wallet kdf")
		}

		// update factor
		if err := q.DeletePasswordFactor(wallet.WalletId); err != nil {
			return errors.Wrap(err, "failed to delete password factor")
		}
		if err := q.CreatePasswordFactor(wallet.WalletId, factor); err != nil {
			return errors.Wrap(err, "failed to create password factor")
		}

		// submit transaction
		// TODO validate envelope not empty
		if result := Horizon(r).Submitter().Submit(r.Context(), request.Data.Relationships.Transaction.Data.Attributes.Envelope); result.Err != nil {
			// TODO assert fail reasons
			return errors.Wrap(result.Err, "failed to submit transaction", result.GetLoganFields())
		}

		return nil
	})
	if err != nil {
		cause := errors.Cause(err)

		if cause == api.ErrWalletsConflict {
			ape.RenderErr(w, problems.Conflict())
			return
		}

		if cause == api.ErrWalletsWalletIDViolated {
			ape.RenderErr(w, problems.Conflict())
			return
		}

		if cause == api.ErrWalletsKDFViolated {
			ape.RenderErr(w, problems.BadRequest(Errors{
				"/data/relationships/kdf/data/id": errors.New("invalid kdf version"),
			})...)
			return
		}

		Log(r).WithError(err).Error("update wallet tx failed")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// render response
	{
		resource := resources.NewWallet(wallet)
		resource.Data.Relationships.Factor = resources.NewPasswordFactor(factor)
		json.NewEncoder(w).Encode(&resource)
	}
}
