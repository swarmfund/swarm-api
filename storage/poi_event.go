package storage

import "gitlab.com/swarmfund/api/db2/api"

func (c *Consumer) ProcessProofOfIncomeUpload(user *api.User, document *api.Document) error {
	document.Meta = map[string]interface{}{
		"reviewed": false,
	}
	err := c.apiQ.Users().Documents(user.DocumentsVersion).Set(user.ID, document)
	if err != nil {
		return err
	}

	return c.apiQ.Users().LimitReviewState(string(user.Address), api.UserLimitReviewPending)
}
