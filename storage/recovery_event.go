package storage

import (
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/utils"
)

func (l *Consumer) ProcessRecoveryUpload(user *api.User, document *api.Document) error {
	recoveryRequest, err := l.apiQ.Recoveries().ByAccountID(string(user.Address))
	if err != nil {
		return err
	}

	if recoveryRequest == nil {
		// TODO delete document, since there is no recovery in progress
		return nil
	}

	if recoveryRequest.RecoveryWalletID != nil {
		// already have doc for this recovery
		// TODO delete document
		return nil
	}

	// recovery document version is recovery wallet ID,
	// but hexed instead of b64. yolo!
	walletID, err := utils.HexToBase64(document.Version)
	if err != nil {
		return errors.Wrap(err, "failed to encode wallet id")
	}
	return l.apiQ.Recoveries().SetRecoveryWalletID(recoveryRequest.ID, walletID)
}
