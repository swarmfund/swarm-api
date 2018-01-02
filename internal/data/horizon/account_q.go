package horizon

import (
	"gitlab.com/swarmfund/go/resources"
)

type AccountQ struct {
	horizon *Client
}

func NewAccountQ(horizon *Client) *AccountQ {
	return &AccountQ{
		horizon: horizon,
	}
}

func (q *AccountQ) Signers(address string) ([]resources.Signer, error) {
	signers, err := q.horizon.Account(address).Signers()
	if err != nil {
		return nil, err
	}

	if signers == nil {
		return nil, nil
	}

	// TODO share resource
	result := make([]resources.Signer, len(signers))
	for _, signer := range signers {
		result = append(result, resources.Signer{
			AccountID:  signer.PublicKey,
			Weight:     int(signer.Weight),
			SignerType: int(signer.SignerTypeI),
			Identity:   int(signer.SignerIdentity),
			Name:       signer.SignerName,
		})
	}
	return result, nil
}
