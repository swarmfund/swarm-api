package storage

import (
	"encoding/json"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/db2/api"
)

func (c *Consumer) ProcessKYCUpload(user *api.User, document *api.Document) error {
	if user.State == api.UserWaitingForApproval {
		return nil
	}

	// TODO block by KYC reason if needed

	// drop document reject reasons
	entity, rr := user.KYCEntities.DocumentsRejectReasons()
	if entity != nil {
		rr.Drop(document.EntityID, document.Type)
		data, err := json.Marshal(&rr)
		if err != nil {
			return errors.Wrap(err, "failed to marshal reject reasons")
		}
		entity.Data = data
		if err = c.apiQ.Users().KYC().Update(entity.ID, entity.Data); err != nil {
			return errors.Wrap(err, "failed to update kyc entity")
		}
	}

	state := user.CheckState()
	if state != user.State {
		if err := c.apiQ.Users().ChangeState(string(user.Address), state); err != nil {
			return errors.Wrap(err, "failed to update state")
		}
	}

	return nil
}
